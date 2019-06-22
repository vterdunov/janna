// repository is a concrete implementations of virtualmachine/usecase.VMRepository interface
package repository

import (
	"archive/tar"
	"context"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"path"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/govc/importx"
	"github.com/vmware/govmomi/nfc"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/ovf"
	"github.com/vmware/govmomi/property"
	"github.com/vmware/govmomi/session"
	"github.com/vmware/govmomi/vim25"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/progress"
	"github.com/vmware/govmomi/vim25/soap"
	"github.com/vmware/govmomi/vim25/types"

	"github.com/vterdunov/janna/internal/virtualmachine/usecase"
)

type VMRepository struct {
	client *govmomi.Client
}

func NewVMRepository(url string, insecure bool) (*VMRepository, error) {
	ctx := context.Background()
	u, err := soap.ParseURL(url)
	if err != nil {
		return nil, err
	}

	soapClient := soap.NewClient(u, insecure)
	vimClient, err := vim25.NewClient(ctx, soapClient)
	if err != nil {
		return nil, err
	}

	vimClient.RoundTripper = session.KeepAlive(vimClient.RoundTripper, 1*time.Minute)
	client := &govmomi.Client{
		Client:         vimClient,
		SessionManager: session.NewManager(vimClient),
	}

	err = client.SessionManager.Login(ctx, u.User)
	if err != nil {
		return nil, err
	}

	repo := VMRepository{
		client: client,
	}

	return &repo, nil
}

func (v *VMRepository) VMInfo(uuid string) (usecase.VMInfoResponse, error) {
	ctx := context.Background()
	vm, err := findByUUID(ctx, v.client.Client, "DC1", uuid)
	if err != nil {
		return usecase.VMInfoResponse{}, err
	}
	refs := make([]types.ManagedObjectReference, 0)
	refs = append(refs, vm.Reference())

	// Retrieve all properties
	// Reference: http://pubs.vmware.com/vsphere-60/topic/com.vmware.wssdk.apiref.doc/vim.VirtualMachine.html
	var mVM mo.VirtualMachine
	var props []string

	pc := property.DefaultCollector(v.client.Client)
	if err := pc.Retrieve(ctx, refs, props, &mVM); err != nil {
		return usecase.VMInfoResponse{}, err
	}

	vmInfo := usecase.VMInfoResponse{
		Name:             mVM.Summary.Config.Name,
		UUID:             mVM.Summary.Config.Uuid,
		Template:         mVM.Summary.Config.Template,
		GuestID:          mVM.Summary.Config.GuestId,
		Annotation:       mVM.Summary.Config.Annotation,
		PowerState:       string(mVM.Runtime.PowerState),
		NumCPU:           uint32(mVM.Summary.Config.NumCpu),
		NumEthernetCards: uint32(mVM.Summary.Config.NumEthernetCards),
		NumVirtualDisks:  uint32(mVM.Summary.Config.NumVirtualDisks),
	}

	return vmInfo, nil
}

func (v *VMRepository) VMDeploy(params usecase.VMDeployRequest) (usecase.VMDeployResponse, error) {
	ctx := context.Background()
	deploy, err := newOVFx(ctx, v.client.Client, params)
	if err != nil {
		return usecase.VMDeployResponse{}, err
	}

	opener := importx.Opener{
		Client: v.client.Client,
	}

	ta := tapeArchive{
		path:   params.OvaURL,
		Opener: opener,
	}

	archive := importx.ArchiveFlag{
		Archive: ta,
	}

	o, err := archive.ReadOvf("*.ovf")
	if err != nil {
		return usecase.VMDeployResponse{}, err
	}

	e, err := archive.ReadEnvelope(o)
	if err != nil {
		return usecase.VMDeployResponse{}, fmt.Errorf("failed to parse ovf: %s", err)
	}

	name := params.Name

	cisp := types.OvfCreateImportSpecParams{
		// See https://github.com/vmware/govmomi/blob/v0.16.0/vim25/types/enum.go#L3381-L3395
		// VMWare can not support some of those disk format types
		// "preallocated", "thin", "seSparse", "rdm", "rdmp",
		// "raw", "delta", "sparse2Gb", "thick2Gb", "eagerZeroedThick",
		// "sparseMonolithic", "flatMonolithic", "thick"
		// TODO: get form params
		DiskProvisioning: "thin",
		EntityName:       name,
		NetworkMapping:   deploy.networkMap(e),
	}

	m := ovf.NewManager(v.client.Client)
	ovfContent := string(o)
	rp := deploy.ResourcePool
	ds := deploy.Datastore

	spec, err := m.CreateImportSpec(ctx, ovfContent, rp, ds, cisp)
	if err != nil {
		return usecase.VMDeployResponse{}, errors.Wrap(err, "Could not create VM spec")
	}
	if spec.Error != nil {
		return usecase.VMDeployResponse{}, errors.New(spec.Error[0].LocalizedMessage)
	}

	if params.Annotation != "" {
		switch s := spec.ImportSpec.(type) {
		case *types.VirtualMachineImportSpec:
			s.ConfigSpec.Annotation = params.Annotation
		case *types.VirtualAppImportSpec:
			s.VAppConfigSpec.Annotation = params.Annotation
		}
	}

	lease, err := rp.ImportVApp(ctx, spec.ImportSpec, deploy.Folder, deploy.Host)
	if err != nil {
		err = errors.Wrap(err, "Could not import Virtual Appliance")
		return usecase.VMDeployResponse{}, err
	}

	info, err := lease.Wait(ctx, spec.FileItem)
	if err != nil {
		err = errors.Wrap(err, "error while waiting lease")
		return usecase.VMDeployResponse{}, err
	}

	u := lease.StartUpdater(ctx, info)
	defer u.Done()
	for _, item := range info.Items {
		if err = deploy.Upload(ctx, lease, item, archive); err != nil {
			return usecase.VMDeployResponse{}, errors.Wrapf(err, "Could not upload disk to VMWare, disk: %v", item.Path)
		}
	}

	return usecase.VMDeployResponse{}, nil
}

type tapeArchive struct {
	path string
	importx.Opener
}

func (t tapeArchive) Open(name string) (io.ReadCloser, int64, error) {
	f, _, err := t.OpenFile(t.path)
	if err != nil {
		return nil, 0, err
	}

	r := tar.NewReader(f)

	for {
		h, err := r.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, 0, err
		}

		matched, err := path.Match(name, path.Base(h.Name))
		if err != nil {
			return nil, 0, err
		}

		if matched {
			return &tapeArchiveEntry{r, f}, h.Size, nil
		}
	}

	_ = f.Close()

	return nil, 0, os.ErrNotExist
}

type tapeArchiveEntry struct {
	io.Reader
	f io.Closer
}

func (t *tapeArchiveEntry) Close() error {
	return t.f.Close()
}

// findByUUID find and returns VM by its UUID
func findByUUID(ctx context.Context, client *vim25.Client, dcName, uuid string) (*object.VirtualMachine, error) {
	f := find.NewFinder(client, true)

	dc, err := f.DatacenterOrDefault(ctx, dcName)
	if err != nil {
		return nil, err
	}

	f.SetDatacenter(dc)

	si := object.NewSearchIndex(client)

	ref, err := si.FindByUuid(ctx, dc, uuid, true, nil)
	if err != nil {
		return nil, err
	}

	vm, ok := ref.(*object.VirtualMachine)
	if !ok {
		return nil, errors.New("could not find Virtual Machine by UUID. Could not assert reference to Virtual Machine")
	}

	return vm, nil
}

type ovfx struct {
	client *vim25.Client
	finder *find.Finder

	Datacenter     *object.Datacenter
	Datastore      *object.Datastore
	Folder         *object.Folder
	ResourcePool   *object.ResourcePool
	Host           *object.HostSystem
	NetworkMapping []importx.Network
}

// newOVFx creates a new (ovfx)deployment object.
// It choose needed resources
func newOVFx(ctx context.Context, client *vim25.Client, params usecase.VMDeployRequest) (*ovfx, error) {
	d := ovfx{
		client: client,
	}

	d.finder = find.NewFinder(client, true)

	// step 1. choose Datacenter and folder
	if err := d.chooseDatacenter(ctx, params.Datacenter); err != nil {
		err = errors.Wrap(err, "Could not choose datacenter")
		return nil, err
	}

	if err := d.chooseFolder(ctx, params.Folder); err != nil {
		err = errors.Wrap(err, "Could not choose folder")
		return nil, err
	}

	// step 2. choose computer resource
	resType := params.ComputerResources.Type
	resPath := params.ComputerResources.Path
	if err := d.chooseComputerResource(ctx, resType, resPath); err != nil {
		err = errors.Wrap(err, "Could not choose Computer Resource")
		return nil, err
	}

	// step 3. Choose datastore cluster or single datastore
	dsType := params.Datastores.Type
	dsNames := params.Datastores.Names
	if err := d.chooseDatastore(ctx, dsType, dsNames); err != nil {
		err = errors.Wrap(err, "Could not choose datastore")
		return nil, err
	}

	return &d, nil
}

func (o *ovfx) networkMap(e *ovf.Envelope) (p []types.OvfNetworkMapping) {
	ctx := context.TODO()
	networks := map[string]string{}

	if e.Network != nil {
		for _, net := range e.Network.Networks {
			networks[net.Name] = net.Name
		}
	}

	for _, net := range o.NetworkMapping {
		networks[net.Name] = net.Network
	}

	for src, dst := range networks {
		if net, err := o.finder.Network(ctx, dst); err == nil {
			p = append(p, types.OvfNetworkMapping{
				Name:    src,
				Network: net.Reference(),
			})
		}
	}

	return p
}

func (o *ovfx) chooseComputerResource(ctx context.Context, resType usecase.ComputerResourcesType, path string) error {
	switch resType {
	case usecase.ComputerResourceHost:
		if err := o.computerResourceWithHost(ctx, path); err != nil {
			return err
		}
	case usecase.ComputerResourceCluster:
		if err := o.computerResourceWithCluster(ctx, path); err != nil {
			return err
		}
	case usecase.ComputerResourceResourcePool:
		if err := o.computerResourceWithResourcePool(ctx, path); err != nil {
			return err
		}
	default:
		return fmt.Errorf("could not recognize computer resource type. type: %q. path: %q", resType, path)
	}

	return nil
}

func (o *ovfx) computerResourceWithHost(ctx context.Context, path string) error {
	host, err := o.finder.HostSystemOrDefault(ctx, path)
	if err != nil {
		return err
	}

	rp, err := host.ResourcePool(ctx)
	if err != nil {
		return err
	}

	o.Host = host
	o.ResourcePool = rp
	return nil
}

func (o *ovfx) computerResourceWithCluster(ctx context.Context, path string) error {
	cluster, err := o.finder.ClusterComputeResourceOrDefault(ctx, path)
	if err != nil {
		return err
	}

	rp, err := cluster.ResourcePool(ctx)
	if err != nil {
		return err
	}

	o.ResourcePool = rp

	// vCenter will choose a host
	o.Host = nil
	return nil
}

func (o *ovfx) computerResourceWithResourcePool(ctx context.Context, rpName string) error {
	rp, err := o.finder.ResourcePoolOrDefault(ctx, rpName)
	if err != nil {
		return err
	}

	o.ResourcePool = rp

	// vCenter will choose a host
	o.Host = nil
	return nil
}

func (o *ovfx) chooseDatacenter(ctx context.Context, dcName string) error {
	dc, err := o.finder.DatacenterOrDefault(ctx, dcName)
	if err != nil {
		return err
	}
	o.finder.SetDatacenter(dc)
	o.Datacenter = dc
	return nil
}

func (o *ovfx) chooseDatastore(ctx context.Context, dsType usecase.DatastoreType, names []string) error {
	switch dsType {
	case usecase.DatastoreCluster:
		if err := o.chooseDatastoreWithCluster(ctx, names); err != nil {
			return err
		}
	case usecase.DatastoreDatastore:
		if err := o.chooseDatastoreWithDatastore(ctx, names); err != nil {
			return err
		}
	default:
		return errors.New("could not recognize datastore type")
	}
	return nil
}

func (o *ovfx) chooseDatastoreWithCluster(ctx context.Context, names []string) error {
	pod, err := o.finder.DatastoreClusterOrDefault(ctx, pickRandom(names))
	if err != nil {
		return err
	}

	drsEnabled, err := isStorageDRSEnabled(ctx, pod)
	if err != nil {
		return err
	}
	if !drsEnabled {
		return errors.New("storage DRS is not enabled on datastore cluster")
	}

	var vmSpec types.VirtualMachineConfigSpec
	sps := types.StoragePlacementSpec{
		Type:         string(types.StoragePlacementSpecPlacementTypeCreate),
		ResourcePool: types.NewReference(o.ResourcePool.Reference()),
		PodSelectionSpec: types.StorageDrsPodSelectionSpec{
			StoragePod: types.NewReference(pod.Reference()),
		},
		Folder:     types.NewReference(o.Folder.Reference()),
		ConfigSpec: &vmSpec,
	}

	srm := object.NewStorageResourceManager(o.client)
	placement, err := srm.RecommendDatastores(ctx, sps)
	if err != nil {
		return err
	}

	recs := placement.Recommendations
	if len(recs) < 1 {
		return errors.New("no storage DRS recommendations were found for the requested action")
	}

	spa, ok := recs[0].Action[0].(*types.StoragePlacementAction)
	if !ok {
		return errors.New("could not get datastore from DRS recomendation")
	}

	ds := spa.Destination
	var mds mo.Datastore
	err = property.DefaultCollector(o.client).RetrieveOne(ctx, ds, []string{"name"}, &mds)
	if err != nil {
		return err
	}

	datastore := object.NewDatastore(o.client, ds)

	o.Datastore = datastore
	return nil
}

func isStorageDRSEnabled(ctx context.Context, pod *object.StoragePod) (bool, error) {
	var props mo.StoragePod
	if err := pod.Properties(ctx, pod.Reference(), nil, &props); err != nil {
		return false, err
	}

	if props.PodStorageDrsEntry == nil {
		return false, nil
	}

	return props.PodStorageDrsEntry.StorageDrsConfig.PodConfig.Enabled, nil
}

func (o *ovfx) chooseDatastoreWithDatastore(ctx context.Context, names []string) error {
	ds, err := o.finder.DatastoreOrDefault(ctx, pickRandom(names))
	if err != nil {
		return err
	}

	o.Datastore = ds
	return nil
}

func pickRandom(slice []string) string {
	rand.Seed(time.Now().Unix())
	return slice[rand.Intn(len(slice))]
}

func (o *ovfx) chooseFolder(ctx context.Context, fName string) error {
	folder, err := o.finder.FolderOrDefault(ctx, fName)
	if err != nil {
		return err
	}
	o.Folder = folder
	return nil
}

func (o *ovfx) Upload(ctx context.Context, lease *nfc.Lease, item nfc.FileItem, archive importx.ArchiveFlag) error {
	file := item.Path

	f, size, err := archive.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	outputStr := path.Base(file)
	logger := log.New(os.Stdout, "pl", 2)
	pl := newProgressLogger(outputStr, logger)
	defer pl.Wait()

	opts := soap.Upload{
		ContentLength: size,
		Progress:      pl,
	}

	return lease.Upload(ctx, item, f, opts)
}

type progressLogger struct {
	prefix string

	wg sync.WaitGroup

	sink   chan chan progress.Report
	done   chan struct{}
	logger *log.Logger
}

func (p *progressLogger) loopA() {
	var err error

	defer p.wg.Done()

	tick := time.NewTicker(5 * time.Second)
	defer tick.Stop()

	called := false

	for stop := false; !stop; {
		select {
		case ch := <-p.sink:
			err = p.loopB(tick, ch)
			stop = true
			called = true
		case <-p.done:
			stop = true
		}
	}

	if err != nil && err != io.EOF {
		// p.logger.Log("err", errors.Wrap(err, "Error with disks uploading"), "file", p.prefix)
		p.logger.Print("Disk uploading error")
	}

	if called {
		// p.logger.Log("msg", "uploaded", "file", p.prefix)
		p.logger.Print("Uploaded")
	}
}

// loopB runs after Sink() has been called.
func (p *progressLogger) loopB(tick *time.Ticker, ch <-chan progress.Report) error {
	var r progress.Report
	var ok bool
	var err error

	for ok = true; ok; {
		select {
		case r, ok = <-ch:
			if !ok {
				break
			}
			err = r.Error()
		case <-tick.C:
			if r != nil {
				pc := fmt.Sprintf("%.0f%%", r.Percentage())
				// p.logger.Log("msg", "uploading disks", "file", p.prefix, "progress", pc)
				p.logger.Printf("progress: %v", pc)
			}
		}
	}

	return err
}

func (p *progressLogger) Wait() {
	close(p.done)
	p.wg.Wait()
}

func (p *progressLogger) Sink() chan<- progress.Report {
	ch := make(chan progress.Report)
	p.sink <- ch
	return ch
}

func newProgressLogger(prefix string, logger *log.Logger) *progressLogger {
	p := &progressLogger{
		prefix: prefix,

		sink:   make(chan chan progress.Report),
		done:   make(chan struct{}),
		logger: logger,
	}

	p.wg.Add(1)

	go p.loopA()

	return p
}

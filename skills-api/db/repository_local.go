package db

import (
	"fmt"
	"math"
	"sync"

	"github.com/imhshekhar47/go-api-core/logger"
	"github.com/imhshekhar47/go-api-core/model"
	"github.com/imhshekhar47/go-api-core/service"
	"github.com/imhshekhar47/go-api-core/utils"
	"github.com/imhshekhar47/hs-taskmaster/skills-api/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var localRepoLogger = logger.GetLogger("db/local")

type RepositoryLocal struct {
	lock     sync.Mutex
	sequence service.SequenceService
	store    *model.Datastore[*pb.Skill]
}

func newRepositoryLocal() *RepositoryLocal {
	return &RepositoryLocal{
		lock:     sync.Mutex{},
		sequence: service.NewDefaultSequenceService(),
		store:    model.NewDatastore[*pb.Skill](100),
	}
}

var defaultInmemoryRpository *RepositoryLocal

func GetRepositoryLocal() *RepositoryLocal {
	if defaultInmemoryRpository == nil {
		defaultInmemoryRpository = newRepositoryLocal()
	}

	demoValues := []string{"testing", "iot", "task"}
	for idx := range demoValues {
		req := &pb.SkillAddRequest{
			Name: demoValues[idx],
		}

		defaultInmemoryRpository.Add(req)
	}

	return defaultInmemoryRpository
}

func (r *RepositoryLocal) List(req *pb.PageRequest) (*pb.SkillQueryResponse, error) {
	localRepoLogger.Traceln("entry: List()")
	index := utils.GetNumOrElse(req.Index, 0)
	limit := utils.GetNumOrElse(req.Limit, 10)
	offset := index * limit

	startIndex := offset + (index * limit)
	if int(startIndex) > int(r.store.Size()) {
		localRepoLogger.Debugln("exit: List(), Error: out of range")
		return nil, status.Errorf(codes.OutOfRange, "page request is out of range")
	}

	endIndex := int(math.Min(float64(startIndex+limit), float64(r.store.Size())))

	if r.store.IsEmpty() {
		localRepoLogger.Debugln("exit: List(), Error: empty")
		return &pb.SkillQueryResponse{
			Range: &pb.PageableResult{
				Records: 0,
				Index:   index,
				Limit:   limit,
			},
			Items: make([]*pb.Skill, 0),
		}, nil
	}

	items, err := r.store.Slice(startIndex, int32(endIndex))
	if err != nil {
		return nil, err
	}

	localRepoLogger.Debugln("exit: List(), ok")
	return &pb.SkillQueryResponse{
		Range: &pb.PageableResult{
			Records: r.store.Size(),
			Index:   index,
			Limit:   limit,
		},
		Items: items,
	}, nil
}

func (r *RepositoryLocal) Add(req *pb.SkillAddRequest) (*pb.Skill, error) {
	localRepoLogger.Debugln("entry: Add()")

	// verify input
	if utils.IsEmpty(req.Name) {
		return nil, status.Errorf(codes.InvalidArgument, "Name is missing")
	}

	newId := fmt.Sprintf("%010d", r.sequence.Next())
	newItem := &pb.Skill{
		Id:   newId,
		Name: req.Name,
	}

	if err := r.store.Add(newItem); err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	localRepoLogger.Debugf("exit: Add(), ok %d", r.store.Size())

	return newItem, nil
}

func (r *RepositoryLocal) Find(req *pb.SkillRequest) (*pb.Skill, error) {
	localRepoLogger.Debugln("entry: Find()")
	item, err := r.store.FindById(req.Id)
	if err != nil {
		localRepoLogger.Debugln("exit: Find(), Error: not found")
		return nil, status.Errorf(codes.NotFound, "NOT FOUND")
	}
	localRepoLogger.Debugln("exit: Find(), ok")
	return *item, nil
}

func (r *RepositoryLocal) Patch(req *pb.Skill) (*pb.Skill, error) {
	localRepoLogger.Traceln("entry: Patch()")
	localRepoLogger.Tracef("req: %s", utils.AsJsonString(req))
	if utils.IsEmpty(string(req.Id)) {
		localRepoLogger.Debugln("exit: Remove(), Error: Missing Id")
		return nil, status.Errorf(codes.InvalidArgument, "field Id is missing")
	}

	err := r.store.Add(req)

	if err != nil {
		localRepoLogger.Traceln("exit: Remove(), Error: Invalid Id")
		return nil, status.Errorf(codes.InvalidArgument, "invalid notification id")
	}

	localRepoLogger.Traceln("exit: Patch(), ok")
	return req, nil
}

func (r *RepositoryLocal) Remove(req *pb.SkillRequest) (*pb.Skill, error) {
	localRepoLogger.Debugln("entry: Remove()")
	localRepoLogger.Tracef("req: %s", utils.AsJsonString(req))
	if utils.IsEmpty(string(req.Id)) {
		localRepoLogger.Debugln("exit: Remove(), Error: Missing Id")
		return nil, status.Errorf(codes.InvalidArgument, "field Id is missing")
	}

	itemRef, err := r.store.Remove(req.Id)

	if err != nil {
		localRepoLogger.Debugln("exit: Remove(), Error: Invalid Id")
		return nil, status.Errorf(codes.InvalidArgument, "invalid notification id")
	}

	localRepoLogger.Debugln("exit: Remove(), ok")
	return *itemRef, nil
}

func (r *RepositoryLocal) Search(*pb.SkillQuery) (*pb.SkillQueryResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method not implemented")
}

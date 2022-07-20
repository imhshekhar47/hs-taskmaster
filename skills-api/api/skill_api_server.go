package api

import (
	"context"

	"github.com/imhshekhar47/go-api-core/logger"
	"github.com/imhshekhar47/hs-taskmaster/skills-api/db"
	"github.com/imhshekhar47/hs-taskmaster/skills-api/pb"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var apiLogger = logger.GetLogger("api/skill")

type SkillApiServer struct {
	pb.UnimplementedSkillServiceServer
	repository db.Repository
}

func NewServer(repository db.Repository) *SkillApiServer {
	return &SkillApiServer{
		repository: repository,
	}
}

func (s *SkillApiServer) Healthcheck(context.Context, *emptypb.Empty) (*pb.Health, error) {
	return &pb.Health{
		Timestamp: timestamppb.Now(),
		Status:    pb.HealthStatus_UP,
	}, nil
}

func (s *SkillApiServer) List(ctx context.Context, req *pb.PageRequest) (*pb.SkillQueryResponse, error) {
	apiLogger.Debugln("entry: List()")
	defer apiLogger.Debugln("exit: List()")
	return s.repository.List(req)
}

func (s *SkillApiServer) Add(ctx context.Context, req *pb.SkillAddRequest) (*pb.Skill, error) {
	apiLogger.Debugln("entry: Add()")
	defer apiLogger.Debugln("exit: Add()")
	return s.repository.Add(req)
}

func (s *SkillApiServer) Get(ctx context.Context, req *pb.SkillRequest) (*pb.Skill, error) {
	apiLogger.Debugln("entry: Get()")
	defer apiLogger.Debugln("exit: Get()")

	return s.repository.Find(req)
}

func (s *SkillApiServer) Patch(ctx context.Context, req *pb.Skill) (*pb.Skill, error) {
	apiLogger.Debugln("entry: Patch()")
	defer apiLogger.Debugln("exit: Patch()")
	return s.repository.Patch(req)
}

func (s *SkillApiServer) Remove(ctx context.Context, req *pb.SkillRequest) (*pb.Skill, error) {
	apiLogger.Debugf("entry: Remove(%s)", req.Id)
	defer apiLogger.Debugf("exit: Remove()")

	return s.repository.Remove(req)
}

func (s *SkillApiServer) Search(ctx context.Context, req *pb.SkillQuery) (*pb.SkillQueryResponse, error) {
	apiLogger.Debugf("entry: Search(%s)", req.Name)
	defer apiLogger.Debugf("exit: Search()")

	return s.repository.Search(req)
}

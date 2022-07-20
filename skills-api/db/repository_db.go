package db

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/imhshekhar47/go-api-core/logger"
	"github.com/imhshekhar47/go-api-core/utils"
	"github.com/imhshekhar47/hs-taskmaster/skills-api/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

var repoLogger = logger.GetLogger("db/local")

type SkillRepository struct {
	db *gorm.DB
}

func NewSkillRepository(db *gorm.DB) *SkillRepository {
	return &SkillRepository{
		db: db,
	}
}

func (r *SkillRepository) List(req *pb.PageRequest) (*pb.SkillQueryResponse, error) {
	repoLogger.Traceln("entry: List()")
	index := utils.GetNumOrElse(req.Index, 0)
	limit := utils.GetNumOrElse(req.Limit, 10)
	offset := index * limit

	var count int64
	r.db.Model(&pb.SkillORM{}).Count(&count)

	var models []pb.SkillORM
	r.db.Offset(int(offset)).Limit(int(limit)).Find(&models)

	skills := make([]*pb.Skill, 0)
	for idx := range models {
		if skill, err := models[idx].ToPB(context.Background()); err == nil {
			skills = append(skills, &skill)
		} else {
			repoLogger.Errorf("error while user model to pb")
		}
	}

	repoLogger.Traceln("exit: List()")
	return &pb.SkillQueryResponse{
		Range: &pb.PageableResult{
			Records: int32(count),
			Index:   index,
			Limit:   limit,
		},
		Items: skills,
	}, nil
}

func (r *SkillRepository) Add(req *pb.SkillAddRequest) (*pb.Skill, error) {
	repoLogger.Traceln("entry: Add()")

	if utils.IsEmpty(req.Name) {
		repoLogger.Traceln("exit: Add(), Error missing Name")
		return nil, status.Errorf(codes.InvalidArgument, "requird field Name is empty")
	}

	skill := pb.Skill{
		Id:   uuid.New().String(),
		Name: req.Name,
	}
	_ = r.db.Transaction(func(tx *gorm.DB) error {
		tx.Create(&skill)
		return nil
	})

	repoLogger.Traceln("exit: Add()")
	return &skill, nil
}

func (r *SkillRepository) Find(req *pb.SkillRequest) (*pb.Skill, error) {
	repoLogger.Traceln("entry: Find()")
	ctx := context.Background()

	skill := &pb.SkillORM{}
	r.db.First(skill, "id=?", req.Id)
	item, err := skill.ToPB(ctx)
	if err != nil {
		repoLogger.Traceln("exit: Find(), protobuf creation")
		return nil, fmt.Errorf("failed to convert to protobuf")
	}

	repoLogger.Traceln("exit: Find()")
	return &item, nil
}

func (r *SkillRepository) Patch(req *pb.Skill) (*pb.Skill, error) {
	repoLogger.Traceln("entry: Patch()")

	skill := &pb.SkillORM{}
	r.db.First(skill, "id=?", req.Id)
	if skill.Id == req.Id {
		skill.Name = req.Name
		r.db.Save(skill)
	} else {
		repoLogger.Tracef("exit: Patch(), skill NOT Found with id '%s'", req.Id)
	}
	item, err := skill.ToPB(context.Background())
	if err != nil {
		repoLogger.Traceln("exit: Patch(), failed to convert to pb")
		return nil, status.Errorf(codes.InvalidArgument, "failed to convert to protobuf")
	}

	repoLogger.Traceln("entry: Patch()")
	return &item, nil
}

func (r *SkillRepository) Remove(req *pb.SkillRequest) (*pb.Skill, error) {
	repoLogger.Tracef("entry: Remove(%s)", req.Id)

	skill := &pb.SkillORM{}
	r.db.Delete(skill, req.Id)
	item, err := skill.ToPB(context.Background())
	if err != nil {
		repoLogger.Traceln("exit: Remove(), failed to convert to pb")
		return nil, status.Errorf(codes.InvalidArgument, "failed to convert to protobuf")
	}

	repoLogger.Traceln("exit: Remove()")
	return &item, nil
}

func (r *SkillRepository) Search(req *pb.SkillQuery) (*pb.SkillQueryResponse, error) {
	repoLogger.Tracef("entry: Search()")

	index := utils.GetNumOrElse(req.Index, 0)
	limit := utils.GetNumOrElse(req.Limit, 10)
	offset := index * limit

	if utils.IsEmpty(req.Name) {
		return nil, status.Errorf(codes.InvalidArgument, "required field Name is empty")
	}

	var models []pb.SkillORM
	r.db.Where("name Like ?", "%"+req.Name+"%").Offset(int(offset)).Limit(int(limit)).Find(models)

	skills := make([]*pb.Skill, len(models))

	for idx := range models {
		item, err := models[idx].ToPB(context.Background())
		if err != nil {
			repoLogger.Errorf("failed to convert ot protobuf")
			return nil, status.Errorf(codes.InvalidArgument, "failed to convert to protobuf")
		}
		skills = append(skills, &item)
	}

	repoLogger.Tracef("exit: Search()")
	return &pb.SkillQueryResponse{
		Range: &pb.PageableResult{
			Records: int32(len(skills)),
			Index:   index,
			Limit:   limit,
		},
		Items: skills,
	}, nil
}

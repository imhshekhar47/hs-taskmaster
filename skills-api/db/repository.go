package db

import "github.com/imhshekhar47/hs-taskmaster/skills-api/pb"

type Repository interface {
	List(*pb.PageRequest) (*pb.SkillQueryResponse, error)
	Add(*pb.SkillAddRequest) (*pb.Skill, error)
	Find(*pb.SkillRequest) (*pb.Skill, error)
	Patch(*pb.Skill) (*pb.Skill, error)
	Remove(*pb.SkillRequest) (*pb.Skill, error)
	Search(*pb.SkillQuery) (*pb.SkillQueryResponse, error)
}

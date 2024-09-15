package models

import "hudson-newey/2rm/src/util"

type Config struct {
	HardDeletePaths []string
	SoftDeletePaths []string
}

func (config Config) ShouldHardDelete(path string) bool {
	isInHardDelete := util.InArray(config.HardDeletePaths, path)
	isInSoftDelete := util.InArray(config.SoftDeletePaths, path)
	return isInHardDelete || isInSoftDelete
}

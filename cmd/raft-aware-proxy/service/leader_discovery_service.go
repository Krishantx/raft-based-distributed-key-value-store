package service

import (
	models "raft-based-kv/internal/models"
)

type Leader_Discovery_Service struct {
	Leader models.Node
}

func (this *Leader_Discovery_Service) getLeaderNode() models.Node {
	return this.Leader
}

func (this *Leader_Discovery_Service) modifyLeader() models.Node {
	this.Leader = models.Node{}
	return this.Leader
}

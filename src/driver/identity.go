package driver

import (
	"context"
	"sync"

	proto "github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/go-kit/kit/log"
	"github.com/golang/protobuf/ptypes/wrappers"
)

type IdentityService struct {
	logger log.Logger

	readyMu sync.RWMutex
	ready   bool
}

func NewIdentityService(logger log.Logger) *IdentityService {
	return &IdentityService{
		logger: logger,
	}
}

func (s *IdentityService) SetReady(ready bool) {
	s.readyMu.Lock()
	s.ready = ready
	s.readyMu.Unlock()
}

func (s *IdentityService) isReady() bool {
	s.readyMu.RLock()
	ready := s.ready
	s.readyMu.RUnlock()
	return ready
}

func (s *IdentityService) GetPluginInfo(context.Context, *proto.GetPluginInfoRequest) (*proto.GetPluginInfoResponse, error) {
	resp := &proto.GetPluginInfoResponse{
		Name:          PluginName,
		VendorVersion: PluginVersion,
	}
	return resp, nil
}

func (s *IdentityService) GetPluginCapabilities(context.Context, *proto.GetPluginCapabilitiesRequest) (*proto.GetPluginCapabilitiesResponse, error) {
	resp := &proto.GetPluginCapabilitiesResponse{
		Capabilities: []*proto.PluginCapability{
			{
				Type: &proto.PluginCapability_Service_{
					Service: &proto.PluginCapability_Service{
						Type: proto.PluginCapability_Service_CONTROLLER_SERVICE,
					},
				},
			},
			{
				Type: &proto.PluginCapability_Service_{
					Service: &proto.PluginCapability_Service{
						Type: proto.PluginCapability_Service_VOLUME_ACCESSIBILITY_CONSTRAINTS,
					},
				},
			},
		},
	}
	return resp, nil
}

func (s *IdentityService) Probe(context.Context, *proto.ProbeRequest) (*proto.ProbeResponse, error) {
	resp := &proto.ProbeResponse{
		Ready: &wrappers.BoolValue{Value: s.isReady()},
	}
	return resp, nil
}

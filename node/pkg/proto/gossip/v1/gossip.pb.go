// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        (unknown)
// source: gossip/v1/gossip.proto

package gossipv1

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type GossipMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Message:
	//	*GossipMessage_SignedObservation
	//	*GossipMessage_SignedHeartbeat
	//	*GossipMessage_SignedVaaWithQuorum
	//	*GossipMessage_SignedObservationRequest
	Message isGossipMessage_Message `protobuf_oneof:"message"`
}

func (x *GossipMessage) Reset() {
	*x = GossipMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gossip_v1_gossip_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GossipMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GossipMessage) ProtoMessage() {}

func (x *GossipMessage) ProtoReflect() protoreflect.Message {
	mi := &file_gossip_v1_gossip_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GossipMessage.ProtoReflect.Descriptor instead.
func (*GossipMessage) Descriptor() ([]byte, []int) {
	return file_gossip_v1_gossip_proto_rawDescGZIP(), []int{0}
}

func (m *GossipMessage) GetMessage() isGossipMessage_Message {
	if m != nil {
		return m.Message
	}
	return nil
}

func (x *GossipMessage) GetSignedObservation() *SignedObservation {
	if x, ok := x.GetMessage().(*GossipMessage_SignedObservation); ok {
		return x.SignedObservation
	}
	return nil
}

func (x *GossipMessage) GetSignedHeartbeat() *SignedHeartbeat {
	if x, ok := x.GetMessage().(*GossipMessage_SignedHeartbeat); ok {
		return x.SignedHeartbeat
	}
	return nil
}

func (x *GossipMessage) GetSignedVaaWithQuorum() *SignedVAAWithQuorum {
	if x, ok := x.GetMessage().(*GossipMessage_SignedVaaWithQuorum); ok {
		return x.SignedVaaWithQuorum
	}
	return nil
}

func (x *GossipMessage) GetSignedObservationRequest() *SignedObservationRequest {
	if x, ok := x.GetMessage().(*GossipMessage_SignedObservationRequest); ok {
		return x.SignedObservationRequest
	}
	return nil
}

type isGossipMessage_Message interface {
	isGossipMessage_Message()
}

type GossipMessage_SignedObservation struct {
	SignedObservation *SignedObservation `protobuf:"bytes,2,opt,name=signed_observation,json=signedObservation,proto3,oneof"`
}

type GossipMessage_SignedHeartbeat struct {
	SignedHeartbeat *SignedHeartbeat `protobuf:"bytes,3,opt,name=signed_heartbeat,json=signedHeartbeat,proto3,oneof"`
}

type GossipMessage_SignedVaaWithQuorum struct {
	SignedVaaWithQuorum *SignedVAAWithQuorum `protobuf:"bytes,4,opt,name=signed_vaa_with_quorum,json=signedVaaWithQuorum,proto3,oneof"`
}

type GossipMessage_SignedObservationRequest struct {
	SignedObservationRequest *SignedObservationRequest `protobuf:"bytes,5,opt,name=signed_observation_request,json=signedObservationRequest,proto3,oneof"`
}

func (*GossipMessage_SignedObservation) isGossipMessage_Message() {}

func (*GossipMessage_SignedHeartbeat) isGossipMessage_Message() {}

func (*GossipMessage_SignedVaaWithQuorum) isGossipMessage_Message() {}

func (*GossipMessage_SignedObservationRequest) isGossipMessage_Message() {}

type SignedHeartbeat struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Serialized Heartbeat message.
	Heartbeat []byte `protobuf:"bytes,1,opt,name=heartbeat,proto3" json:"heartbeat,omitempty"`
	// ECDSA signature using the node's guardian public key.
	Signature []byte `protobuf:"bytes,2,opt,name=signature,proto3" json:"signature,omitempty"`
	// Guardian address that signed this payload (truncated Eth address).
	// This is already contained in Heartbeat, however, we want to verify
	// the payload before we deserialize it.
	GuardianAddr []byte `protobuf:"bytes,3,opt,name=guardian_addr,json=guardianAddr,proto3" json:"guardian_addr,omitempty"`
}

func (x *SignedHeartbeat) Reset() {
	*x = SignedHeartbeat{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gossip_v1_gossip_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SignedHeartbeat) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SignedHeartbeat) ProtoMessage() {}

func (x *SignedHeartbeat) ProtoReflect() protoreflect.Message {
	mi := &file_gossip_v1_gossip_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SignedHeartbeat.ProtoReflect.Descriptor instead.
func (*SignedHeartbeat) Descriptor() ([]byte, []int) {
	return file_gossip_v1_gossip_proto_rawDescGZIP(), []int{1}
}

func (x *SignedHeartbeat) GetHeartbeat() []byte {
	if x != nil {
		return x.Heartbeat
	}
	return nil
}

func (x *SignedHeartbeat) GetSignature() []byte {
	if x != nil {
		return x.Signature
	}
	return nil
}

func (x *SignedHeartbeat) GetGuardianAddr() []byte {
	if x != nil {
		return x.GuardianAddr
	}
	return nil
}

// P2P gossip heartbeats for network introspection purposes.
type Heartbeat struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The node's arbitrarily chosen, untrusted nodeName.
	NodeName string `protobuf:"bytes,1,opt,name=node_name,json=nodeName,proto3" json:"node_name,omitempty"`
	// A monotonic counter that resets to zero on startup.
	Counter int64 `protobuf:"varint,2,opt,name=counter,proto3" json:"counter,omitempty"`
	// UNIX wall time.
	Timestamp int64                `protobuf:"varint,3,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	Networks  []*Heartbeat_Network `protobuf:"bytes,4,rep,name=networks,proto3" json:"networks,omitempty"`
	// Human-readable representation of the current bridge node release.
	Version string `protobuf:"bytes,5,opt,name=version,proto3" json:"version,omitempty"`
	// Human-readable representation of the guardian key's address.
	GuardianAddr string `protobuf:"bytes,6,opt,name=guardian_addr,json=guardianAddr,proto3" json:"guardian_addr,omitempty"`
	// UNIX boot timestamp.
	BootTimestamp int64 `protobuf:"varint,7,opt,name=boot_timestamp,json=bootTimestamp,proto3" json:"boot_timestamp,omitempty"`
	// List of features enabled on this node.
	Features []string `protobuf:"bytes,8,rep,name=features,proto3" json:"features,omitempty"`
}

func (x *Heartbeat) Reset() {
	*x = Heartbeat{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gossip_v1_gossip_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Heartbeat) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Heartbeat) ProtoMessage() {}

func (x *Heartbeat) ProtoReflect() protoreflect.Message {
	mi := &file_gossip_v1_gossip_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Heartbeat.ProtoReflect.Descriptor instead.
func (*Heartbeat) Descriptor() ([]byte, []int) {
	return file_gossip_v1_gossip_proto_rawDescGZIP(), []int{2}
}

func (x *Heartbeat) GetNodeName() string {
	if x != nil {
		return x.NodeName
	}
	return ""
}

func (x *Heartbeat) GetCounter() int64 {
	if x != nil {
		return x.Counter
	}
	return 0
}

func (x *Heartbeat) GetTimestamp() int64 {
	if x != nil {
		return x.Timestamp
	}
	return 0
}

func (x *Heartbeat) GetNetworks() []*Heartbeat_Network {
	if x != nil {
		return x.Networks
	}
	return nil
}

func (x *Heartbeat) GetVersion() string {
	if x != nil {
		return x.Version
	}
	return ""
}

func (x *Heartbeat) GetGuardianAddr() string {
	if x != nil {
		return x.GuardianAddr
	}
	return ""
}

func (x *Heartbeat) GetBootTimestamp() int64 {
	if x != nil {
		return x.BootTimestamp
	}
	return 0
}

func (x *Heartbeat) GetFeatures() []string {
	if x != nil {
		return x.Features
	}
	return nil
}

// A SignedObservation is a signed statement by a given guardian node
// that they observed a given event.
//
// Observations always result from an external, final event being observed.
// Examples are emitted messages in finalized blocks on a block or guardian set changes
// injected by node operators after reaching off-chain consensus.
//
// The event is uniquely identified by its hashed (tx_hash, nonce, values...) tuple.
//
// Other nodes will verify the signature. Once any node has observed a quorum of
// guardians submitting valid signatures for a given hash, they can be assembled into a VAA.
//
// Messages without valid signature are dropped unceremoniously.
type SignedObservation struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Guardian pubkey as truncated eth address.
	Addr []byte `protobuf:"bytes,1,opt,name=addr,proto3" json:"addr,omitempty"`
	// The observation's deterministic, unique hash.
	Hash []byte `protobuf:"bytes,2,opt,name=hash,proto3" json:"hash,omitempty"`
	// ECSDA signature of the hash using the node's guardian key.
	Signature []byte `protobuf:"bytes,3,opt,name=signature,proto3" json:"signature,omitempty"`
	// Transaction hash this observation was made from.
	// Optional, included for observability.
	TxHash []byte `protobuf:"bytes,4,opt,name=tx_hash,json=txHash,proto3" json:"tx_hash,omitempty"`
	// Message ID (chain/emitter/seq) for this observation.
	// Optional, included for observability.
	MessageId string `protobuf:"bytes,5,opt,name=message_id,json=messageId,proto3" json:"message_id,omitempty"`
}

func (x *SignedObservation) Reset() {
	*x = SignedObservation{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gossip_v1_gossip_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SignedObservation) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SignedObservation) ProtoMessage() {}

func (x *SignedObservation) ProtoReflect() protoreflect.Message {
	mi := &file_gossip_v1_gossip_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SignedObservation.ProtoReflect.Descriptor instead.
func (*SignedObservation) Descriptor() ([]byte, []int) {
	return file_gossip_v1_gossip_proto_rawDescGZIP(), []int{3}
}

func (x *SignedObservation) GetAddr() []byte {
	if x != nil {
		return x.Addr
	}
	return nil
}

func (x *SignedObservation) GetHash() []byte {
	if x != nil {
		return x.Hash
	}
	return nil
}

func (x *SignedObservation) GetSignature() []byte {
	if x != nil {
		return x.Signature
	}
	return nil
}

func (x *SignedObservation) GetTxHash() []byte {
	if x != nil {
		return x.TxHash
	}
	return nil
}

func (x *SignedObservation) GetMessageId() string {
	if x != nil {
		return x.MessageId
	}
	return ""
}

// A SignedVAAWithQuorum message is sent by nodes whenever one of the VAAs they observed
// reached a 2/3+ quorum to be considered valid. Signed VAAs are broadcasted to the gossip
// network to allow nodes to persist them even if they failed to observe the signature.
type SignedVAAWithQuorum struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Vaa []byte `protobuf:"bytes,1,opt,name=vaa,proto3" json:"vaa,omitempty"`
}

func (x *SignedVAAWithQuorum) Reset() {
	*x = SignedVAAWithQuorum{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gossip_v1_gossip_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SignedVAAWithQuorum) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SignedVAAWithQuorum) ProtoMessage() {}

func (x *SignedVAAWithQuorum) ProtoReflect() protoreflect.Message {
	mi := &file_gossip_v1_gossip_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SignedVAAWithQuorum.ProtoReflect.Descriptor instead.
func (*SignedVAAWithQuorum) Descriptor() ([]byte, []int) {
	return file_gossip_v1_gossip_proto_rawDescGZIP(), []int{4}
}

func (x *SignedVAAWithQuorum) GetVaa() []byte {
	if x != nil {
		return x.Vaa
	}
	return nil
}

// Any guardian can send a SignedObservationRequest to the network to request
// all guardians to re-observe a given transaction. This is rate-limited to one
// request per second per guardian to prevent abuse.
//
// In the current implementation, this is only implemented for Solana.
// For Solana, the tx_hash is the account address of the transaction's message account.
type SignedObservationRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Serialized observation request.
	ObservationRequest []byte `protobuf:"bytes,1,opt,name=observation_request,json=observationRequest,proto3" json:"observation_request,omitempty"`
	// Signature
	Signature    []byte `protobuf:"bytes,2,opt,name=signature,proto3" json:"signature,omitempty"`
	GuardianAddr []byte `protobuf:"bytes,3,opt,name=guardian_addr,json=guardianAddr,proto3" json:"guardian_addr,omitempty"`
}

func (x *SignedObservationRequest) Reset() {
	*x = SignedObservationRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gossip_v1_gossip_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SignedObservationRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SignedObservationRequest) ProtoMessage() {}

func (x *SignedObservationRequest) ProtoReflect() protoreflect.Message {
	mi := &file_gossip_v1_gossip_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SignedObservationRequest.ProtoReflect.Descriptor instead.
func (*SignedObservationRequest) Descriptor() ([]byte, []int) {
	return file_gossip_v1_gossip_proto_rawDescGZIP(), []int{5}
}

func (x *SignedObservationRequest) GetObservationRequest() []byte {
	if x != nil {
		return x.ObservationRequest
	}
	return nil
}

func (x *SignedObservationRequest) GetSignature() []byte {
	if x != nil {
		return x.Signature
	}
	return nil
}

func (x *SignedObservationRequest) GetGuardianAddr() []byte {
	if x != nil {
		return x.GuardianAddr
	}
	return nil
}

type ObservationRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ChainId uint32 `protobuf:"varint,1,opt,name=chain_id,json=chainId,proto3" json:"chain_id,omitempty"`
	TxHash  []byte `protobuf:"bytes,2,opt,name=tx_hash,json=txHash,proto3" json:"tx_hash,omitempty"`
}

func (x *ObservationRequest) Reset() {
	*x = ObservationRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gossip_v1_gossip_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ObservationRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ObservationRequest) ProtoMessage() {}

func (x *ObservationRequest) ProtoReflect() protoreflect.Message {
	mi := &file_gossip_v1_gossip_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ObservationRequest.ProtoReflect.Descriptor instead.
func (*ObservationRequest) Descriptor() ([]byte, []int) {
	return file_gossip_v1_gossip_proto_rawDescGZIP(), []int{6}
}

func (x *ObservationRequest) GetChainId() uint32 {
	if x != nil {
		return x.ChainId
	}
	return 0
}

func (x *ObservationRequest) GetTxHash() []byte {
	if x != nil {
		return x.TxHash
	}
	return nil
}

type Heartbeat_Network struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Canonical chain ID.
	Id uint32 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	// Consensus height of the node.
	Height int64 `protobuf:"varint,2,opt,name=height,proto3" json:"height,omitempty"`
	// Chain-specific human-readable representation of the bridge contract address.
	ContractAddress string `protobuf:"bytes,3,opt,name=contract_address,json=contractAddress,proto3" json:"contract_address,omitempty"`
	// Connection error count
	ErrorCount uint64 `protobuf:"varint,4,opt,name=error_count,json=errorCount,proto3" json:"error_count,omitempty"`
}

func (x *Heartbeat_Network) Reset() {
	*x = Heartbeat_Network{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gossip_v1_gossip_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Heartbeat_Network) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Heartbeat_Network) ProtoMessage() {}

func (x *Heartbeat_Network) ProtoReflect() protoreflect.Message {
	mi := &file_gossip_v1_gossip_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Heartbeat_Network.ProtoReflect.Descriptor instead.
func (*Heartbeat_Network) Descriptor() ([]byte, []int) {
	return file_gossip_v1_gossip_proto_rawDescGZIP(), []int{2, 0}
}

func (x *Heartbeat_Network) GetId() uint32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Heartbeat_Network) GetHeight() int64 {
	if x != nil {
		return x.Height
	}
	return 0
}

func (x *Heartbeat_Network) GetContractAddress() string {
	if x != nil {
		return x.ContractAddress
	}
	return ""
}

func (x *Heartbeat_Network) GetErrorCount() uint64 {
	if x != nil {
		return x.ErrorCount
	}
	return 0
}

var File_gossip_v1_gossip_proto protoreflect.FileDescriptor

var file_gossip_v1_gossip_proto_rawDesc = []byte{
	0x0a, 0x16, 0x67, 0x6f, 0x73, 0x73, 0x69, 0x70, 0x2f, 0x76, 0x31, 0x2f, 0x67, 0x6f, 0x73, 0x73,
	0x69, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x09, 0x67, 0x6f, 0x73, 0x73, 0x69, 0x70,
	0x2e, 0x76, 0x31, 0x22, 0xee, 0x02, 0x0a, 0x0d, 0x47, 0x6f, 0x73, 0x73, 0x69, 0x70, 0x4d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x4d, 0x0a, 0x12, 0x73, 0x69, 0x67, 0x6e, 0x65, 0x64, 0x5f,
	0x6f, 0x62, 0x73, 0x65, 0x72, 0x76, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x1c, 0x2e, 0x67, 0x6f, 0x73, 0x73, 0x69, 0x70, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x69,
	0x67, 0x6e, 0x65, 0x64, 0x4f, 0x62, 0x73, 0x65, 0x72, 0x76, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x48,
	0x00, 0x52, 0x11, 0x73, 0x69, 0x67, 0x6e, 0x65, 0x64, 0x4f, 0x62, 0x73, 0x65, 0x72, 0x76, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x12, 0x47, 0x0a, 0x10, 0x73, 0x69, 0x67, 0x6e, 0x65, 0x64, 0x5f, 0x68,
	0x65, 0x61, 0x72, 0x74, 0x62, 0x65, 0x61, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a,
	0x2e, 0x67, 0x6f, 0x73, 0x73, 0x69, 0x70, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x69, 0x67, 0x6e, 0x65,
	0x64, 0x48, 0x65, 0x61, 0x72, 0x74, 0x62, 0x65, 0x61, 0x74, 0x48, 0x00, 0x52, 0x0f, 0x73, 0x69,
	0x67, 0x6e, 0x65, 0x64, 0x48, 0x65, 0x61, 0x72, 0x74, 0x62, 0x65, 0x61, 0x74, 0x12, 0x55, 0x0a,
	0x16, 0x73, 0x69, 0x67, 0x6e, 0x65, 0x64, 0x5f, 0x76, 0x61, 0x61, 0x5f, 0x77, 0x69, 0x74, 0x68,
	0x5f, 0x71, 0x75, 0x6f, 0x72, 0x75, 0x6d, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e,
	0x67, 0x6f, 0x73, 0x73, 0x69, 0x70, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x69, 0x67, 0x6e, 0x65, 0x64,
	0x56, 0x41, 0x41, 0x57, 0x69, 0x74, 0x68, 0x51, 0x75, 0x6f, 0x72, 0x75, 0x6d, 0x48, 0x00, 0x52,
	0x13, 0x73, 0x69, 0x67, 0x6e, 0x65, 0x64, 0x56, 0x61, 0x61, 0x57, 0x69, 0x74, 0x68, 0x51, 0x75,
	0x6f, 0x72, 0x75, 0x6d, 0x12, 0x63, 0x0a, 0x1a, 0x73, 0x69, 0x67, 0x6e, 0x65, 0x64, 0x5f, 0x6f,
	0x62, 0x73, 0x65, 0x72, 0x76, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x72, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x23, 0x2e, 0x67, 0x6f, 0x73, 0x73, 0x69,
	0x70, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x69, 0x67, 0x6e, 0x65, 0x64, 0x4f, 0x62, 0x73, 0x65, 0x72,
	0x76, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x48, 0x00, 0x52,
	0x18, 0x73, 0x69, 0x67, 0x6e, 0x65, 0x64, 0x4f, 0x62, 0x73, 0x65, 0x72, 0x76, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x42, 0x09, 0x0a, 0x07, 0x6d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x22, 0x72, 0x0a, 0x0f, 0x53, 0x69, 0x67, 0x6e, 0x65, 0x64, 0x48, 0x65,
	0x61, 0x72, 0x74, 0x62, 0x65, 0x61, 0x74, 0x12, 0x1c, 0x0a, 0x09, 0x68, 0x65, 0x61, 0x72, 0x74,
	0x62, 0x65, 0x61, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x09, 0x68, 0x65, 0x61, 0x72,
	0x74, 0x62, 0x65, 0x61, 0x74, 0x12, 0x1c, 0x0a, 0x09, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75,
	0x72, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x09, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x74,
	0x75, 0x72, 0x65, 0x12, 0x23, 0x0a, 0x0d, 0x67, 0x75, 0x61, 0x72, 0x64, 0x69, 0x61, 0x6e, 0x5f,
	0x61, 0x64, 0x64, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0c, 0x67, 0x75, 0x61, 0x72,
	0x64, 0x69, 0x61, 0x6e, 0x41, 0x64, 0x64, 0x72, 0x22, 0x9b, 0x03, 0x0a, 0x09, 0x48, 0x65, 0x61,
	0x72, 0x74, 0x62, 0x65, 0x61, 0x74, 0x12, 0x1b, 0x0a, 0x09, 0x6e, 0x6f, 0x64, 0x65, 0x5f, 0x6e,
	0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6e, 0x6f, 0x64, 0x65, 0x4e,
	0x61, 0x6d, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x65, 0x72, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x07, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x65, 0x72, 0x12, 0x1c, 0x0a,
	0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x12, 0x38, 0x0a, 0x08, 0x6e,
	0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1c, 0x2e,
	0x67, 0x6f, 0x73, 0x73, 0x69, 0x70, 0x2e, 0x76, 0x31, 0x2e, 0x48, 0x65, 0x61, 0x72, 0x74, 0x62,
	0x65, 0x61, 0x74, 0x2e, 0x4e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x52, 0x08, 0x6e, 0x65, 0x74,
	0x77, 0x6f, 0x72, 0x6b, 0x73, 0x12, 0x18, 0x0a, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e,
	0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12,
	0x23, 0x0a, 0x0d, 0x67, 0x75, 0x61, 0x72, 0x64, 0x69, 0x61, 0x6e, 0x5f, 0x61, 0x64, 0x64, 0x72,
	0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x67, 0x75, 0x61, 0x72, 0x64, 0x69, 0x61, 0x6e,
	0x41, 0x64, 0x64, 0x72, 0x12, 0x25, 0x0a, 0x0e, 0x62, 0x6f, 0x6f, 0x74, 0x5f, 0x74, 0x69, 0x6d,
	0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x07, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0d, 0x62, 0x6f,
	0x6f, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x12, 0x1a, 0x0a, 0x08, 0x66,
	0x65, 0x61, 0x74, 0x75, 0x72, 0x65, 0x73, 0x18, 0x08, 0x20, 0x03, 0x28, 0x09, 0x52, 0x08, 0x66,
	0x65, 0x61, 0x74, 0x75, 0x72, 0x65, 0x73, 0x1a, 0x7d, 0x0a, 0x07, 0x4e, 0x65, 0x74, 0x77, 0x6f,
	0x72, 0x6b, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x02,
	0x69, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x68, 0x65, 0x69, 0x67, 0x68, 0x74, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x06, 0x68, 0x65, 0x69, 0x67, 0x68, 0x74, 0x12, 0x29, 0x0a, 0x10, 0x63, 0x6f,
	0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x5f, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0f, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x41, 0x64,
	0x64, 0x72, 0x65, 0x73, 0x73, 0x12, 0x1f, 0x0a, 0x0b, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x5f, 0x63,
	0x6f, 0x75, 0x6e, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0a, 0x65, 0x72, 0x72, 0x6f,
	0x72, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x22, 0x91, 0x01, 0x0a, 0x11, 0x53, 0x69, 0x67, 0x6e, 0x65,
	0x64, 0x4f, 0x62, 0x73, 0x65, 0x72, 0x76, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x12, 0x0a, 0x04,
	0x61, 0x64, 0x64, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x61, 0x64, 0x64, 0x72,
	0x12, 0x12, 0x0a, 0x04, 0x68, 0x61, 0x73, 0x68, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04,
	0x68, 0x61, 0x73, 0x68, 0x12, 0x1c, 0x0a, 0x09, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72,
	0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x09, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75,
	0x72, 0x65, 0x12, 0x17, 0x0a, 0x07, 0x74, 0x78, 0x5f, 0x68, 0x61, 0x73, 0x68, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x0c, 0x52, 0x06, 0x74, 0x78, 0x48, 0x61, 0x73, 0x68, 0x12, 0x1d, 0x0a, 0x0a, 0x6d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x09, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x49, 0x64, 0x22, 0x27, 0x0a, 0x13, 0x53, 0x69,
	0x67, 0x6e, 0x65, 0x64, 0x56, 0x41, 0x41, 0x57, 0x69, 0x74, 0x68, 0x51, 0x75, 0x6f, 0x72, 0x75,
	0x6d, 0x12, 0x10, 0x0a, 0x03, 0x76, 0x61, 0x61, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x03,
	0x76, 0x61, 0x61, 0x22, 0x8e, 0x01, 0x0a, 0x18, 0x53, 0x69, 0x67, 0x6e, 0x65, 0x64, 0x4f, 0x62,
	0x73, 0x65, 0x72, 0x76, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x2f, 0x0a, 0x13, 0x6f, 0x62, 0x73, 0x65, 0x72, 0x76, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f,
	0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x12, 0x6f,
	0x62, 0x73, 0x65, 0x72, 0x76, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x1c, 0x0a, 0x09, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0c, 0x52, 0x09, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x12,
	0x23, 0x0a, 0x0d, 0x67, 0x75, 0x61, 0x72, 0x64, 0x69, 0x61, 0x6e, 0x5f, 0x61, 0x64, 0x64, 0x72,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0c, 0x67, 0x75, 0x61, 0x72, 0x64, 0x69, 0x61, 0x6e,
	0x41, 0x64, 0x64, 0x72, 0x22, 0x48, 0x0a, 0x12, 0x4f, 0x62, 0x73, 0x65, 0x72, 0x76, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x19, 0x0a, 0x08, 0x63, 0x68,
	0x61, 0x69, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x07, 0x63, 0x68,
	0x61, 0x69, 0x6e, 0x49, 0x64, 0x12, 0x17, 0x0a, 0x07, 0x74, 0x78, 0x5f, 0x68, 0x61, 0x73, 0x68,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x06, 0x74, 0x78, 0x48, 0x61, 0x73, 0x68, 0x42, 0x41,
	0x5a, 0x3f, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x63, 0x65, 0x72,
	0x74, 0x75, 0x73, 0x6f, 0x6e, 0x65, 0x2f, 0x77, 0x6f, 0x72, 0x6d, 0x68, 0x6f, 0x6c, 0x65, 0x2f,
	0x6e, 0x6f, 0x64, 0x65, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x67,
	0x6f, 0x73, 0x73, 0x69, 0x70, 0x2f, 0x76, 0x31, 0x3b, 0x67, 0x6f, 0x73, 0x73, 0x69, 0x70, 0x76,
	0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_gossip_v1_gossip_proto_rawDescOnce sync.Once
	file_gossip_v1_gossip_proto_rawDescData = file_gossip_v1_gossip_proto_rawDesc
)

func file_gossip_v1_gossip_proto_rawDescGZIP() []byte {
	file_gossip_v1_gossip_proto_rawDescOnce.Do(func() {
		file_gossip_v1_gossip_proto_rawDescData = protoimpl.X.CompressGZIP(file_gossip_v1_gossip_proto_rawDescData)
	})
	return file_gossip_v1_gossip_proto_rawDescData
}

var file_gossip_v1_gossip_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_gossip_v1_gossip_proto_goTypes = []interface{}{
	(*GossipMessage)(nil),            // 0: gossip.v1.GossipMessage
	(*SignedHeartbeat)(nil),          // 1: gossip.v1.SignedHeartbeat
	(*Heartbeat)(nil),                // 2: gossip.v1.Heartbeat
	(*SignedObservation)(nil),        // 3: gossip.v1.SignedObservation
	(*SignedVAAWithQuorum)(nil),      // 4: gossip.v1.SignedVAAWithQuorum
	(*SignedObservationRequest)(nil), // 5: gossip.v1.SignedObservationRequest
	(*ObservationRequest)(nil),       // 6: gossip.v1.ObservationRequest
	(*Heartbeat_Network)(nil),        // 7: gossip.v1.Heartbeat.Network
}
var file_gossip_v1_gossip_proto_depIdxs = []int32{
	3, // 0: gossip.v1.GossipMessage.signed_observation:type_name -> gossip.v1.SignedObservation
	1, // 1: gossip.v1.GossipMessage.signed_heartbeat:type_name -> gossip.v1.SignedHeartbeat
	4, // 2: gossip.v1.GossipMessage.signed_vaa_with_quorum:type_name -> gossip.v1.SignedVAAWithQuorum
	5, // 3: gossip.v1.GossipMessage.signed_observation_request:type_name -> gossip.v1.SignedObservationRequest
	7, // 4: gossip.v1.Heartbeat.networks:type_name -> gossip.v1.Heartbeat.Network
	5, // [5:5] is the sub-list for method output_type
	5, // [5:5] is the sub-list for method input_type
	5, // [5:5] is the sub-list for extension type_name
	5, // [5:5] is the sub-list for extension extendee
	0, // [0:5] is the sub-list for field type_name
}

func init() { file_gossip_v1_gossip_proto_init() }
func file_gossip_v1_gossip_proto_init() {
	if File_gossip_v1_gossip_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_gossip_v1_gossip_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GossipMessage); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_gossip_v1_gossip_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SignedHeartbeat); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_gossip_v1_gossip_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Heartbeat); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_gossip_v1_gossip_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SignedObservation); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_gossip_v1_gossip_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SignedVAAWithQuorum); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_gossip_v1_gossip_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SignedObservationRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_gossip_v1_gossip_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ObservationRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_gossip_v1_gossip_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Heartbeat_Network); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	file_gossip_v1_gossip_proto_msgTypes[0].OneofWrappers = []interface{}{
		(*GossipMessage_SignedObservation)(nil),
		(*GossipMessage_SignedHeartbeat)(nil),
		(*GossipMessage_SignedVaaWithQuorum)(nil),
		(*GossipMessage_SignedObservationRequest)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_gossip_v1_gossip_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_gossip_v1_gossip_proto_goTypes,
		DependencyIndexes: file_gossip_v1_gossip_proto_depIdxs,
		MessageInfos:      file_gossip_v1_gossip_proto_msgTypes,
	}.Build()
	File_gossip_v1_gossip_proto = out.File
	file_gossip_v1_gossip_proto_rawDesc = nil
	file_gossip_v1_gossip_proto_goTypes = nil
	file_gossip_v1_gossip_proto_depIdxs = nil
}

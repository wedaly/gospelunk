// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.19.1
// source: db.proto

package protobuf

import (
	reflect "reflect"
	sync "sync"

	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Package struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name       string    `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Dir        string    `protobuf:"bytes,2,opt,name=dir,proto3" json:"dir,omitempty"`
	ImportPath string    `protobuf:"bytes,3,opt,name=import_path,json=importPath,proto3" json:"import_path,omitempty"`
	GoFiles    []*GoFile `protobuf:"bytes,4,rep,name=go_files,json=goFiles,proto3" json:"go_files,omitempty"`
	Imports    []string  `protobuf:"bytes,5,rep,name=imports,proto3" json:"imports,omitempty"`
	Hash       string    `protobuf:"bytes,6,opt,name=hash,proto3" json:"hash,omitempty"`
}

func (x *Package) Reset() {
	*x = Package{}
	if protoimpl.UnsafeEnabled {
		mi := &file_db_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Package) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Package) ProtoMessage() {}

func (x *Package) ProtoReflect() protoreflect.Message {
	mi := &file_db_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Package.ProtoReflect.Descriptor instead.
func (*Package) Descriptor() ([]byte, []int) {
	return file_db_proto_rawDescGZIP(), []int{0}
}

func (x *Package) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Package) GetDir() string {
	if x != nil {
		return x.Dir
	}
	return ""
}

func (x *Package) GetImportPath() string {
	if x != nil {
		return x.ImportPath
	}
	return ""
}

func (x *Package) GetGoFiles() []*GoFile {
	if x != nil {
		return x.GoFiles
	}
	return nil
}

func (x *Package) GetImports() []string {
	if x != nil {
		return x.Imports
	}
	return nil
}

func (x *Package) GetHash() string {
	if x != nil {
		return x.Hash
	}
	return ""
}

type GoFile struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Filename string   `protobuf:"bytes,1,opt,name=filename,proto3" json:"filename,omitempty"`
	Defs     []*GoDef `protobuf:"bytes,2,rep,name=defs,proto3" json:"defs,omitempty"`
}

func (x *GoFile) Reset() {
	*x = GoFile{}
	if protoimpl.UnsafeEnabled {
		mi := &file_db_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GoFile) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GoFile) ProtoMessage() {}

func (x *GoFile) ProtoReflect() protoreflect.Message {
	mi := &file_db_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GoFile.ProtoReflect.Descriptor instead.
func (*GoFile) Descriptor() ([]byte, []int) {
	return file_db_proto_rawDescGZIP(), []int{1}
}

func (x *GoFile) GetFilename() string {
	if x != nil {
		return x.Filename
	}
	return ""
}

func (x *GoFile) GetDefs() []*GoDef {
	if x != nil {
		return x.Defs
	}
	return nil
}

type GoDef struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name     string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Kind     int32  `protobuf:"varint,2,opt,name=kind,proto3" json:"kind,omitempty"`
	LineNum  int64  `protobuf:"varint,3,opt,name=line_num,json=lineNum,proto3" json:"line_num,omitempty"`
	Exported bool   `protobuf:"varint,4,opt,name=exported,proto3" json:"exported,omitempty"`
}

func (x *GoDef) Reset() {
	*x = GoDef{}
	if protoimpl.UnsafeEnabled {
		mi := &file_db_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GoDef) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GoDef) ProtoMessage() {}

func (x *GoDef) ProtoReflect() protoreflect.Message {
	mi := &file_db_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GoDef.ProtoReflect.Descriptor instead.
func (*GoDef) Descriptor() ([]byte, []int) {
	return file_db_proto_rawDescGZIP(), []int{2}
}

func (x *GoDef) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *GoDef) GetKind() int32 {
	if x != nil {
		return x.Kind
	}
	return 0
}

func (x *GoDef) GetLineNum() int64 {
	if x != nil {
		return x.LineNum
	}
	return 0
}

func (x *GoDef) GetExported() bool {
	if x != nil {
		return x.Exported
	}
	return false
}

var File_db_proto protoreflect.FileDescriptor

var file_db_proto_rawDesc = []byte{
	0x0a, 0x08, 0x64, 0x62, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xa2, 0x01, 0x0a, 0x07, 0x50,
	0x61, 0x63, 0x6b, 0x61, 0x67, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x64, 0x69,
	0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x64, 0x69, 0x72, 0x12, 0x1f, 0x0a, 0x0b,
	0x69, 0x6d, 0x70, 0x6f, 0x72, 0x74, 0x5f, 0x70, 0x61, 0x74, 0x68, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0a, 0x69, 0x6d, 0x70, 0x6f, 0x72, 0x74, 0x50, 0x61, 0x74, 0x68, 0x12, 0x22, 0x0a,
	0x08, 0x67, 0x6f, 0x5f, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x07, 0x2e, 0x47, 0x6f, 0x46, 0x69, 0x6c, 0x65, 0x52, 0x07, 0x67, 0x6f, 0x46, 0x69, 0x6c, 0x65,
	0x73, 0x12, 0x18, 0x0a, 0x07, 0x69, 0x6d, 0x70, 0x6f, 0x72, 0x74, 0x73, 0x18, 0x05, 0x20, 0x03,
	0x28, 0x09, 0x52, 0x07, 0x69, 0x6d, 0x70, 0x6f, 0x72, 0x74, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x68,
	0x61, 0x73, 0x68, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x68, 0x61, 0x73, 0x68, 0x22,
	0x40, 0x0a, 0x06, 0x47, 0x6f, 0x46, 0x69, 0x6c, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x66, 0x69, 0x6c,
	0x65, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x66, 0x69, 0x6c,
	0x65, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1a, 0x0a, 0x04, 0x64, 0x65, 0x66, 0x73, 0x18, 0x02, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x06, 0x2e, 0x47, 0x6f, 0x44, 0x65, 0x66, 0x52, 0x04, 0x64, 0x65, 0x66,
	0x73, 0x22, 0x66, 0x0a, 0x05, 0x47, 0x6f, 0x44, 0x65, 0x66, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61,
	0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x12,
	0x0a, 0x04, 0x6b, 0x69, 0x6e, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x6b, 0x69,
	0x6e, 0x64, 0x12, 0x19, 0x0a, 0x08, 0x6c, 0x69, 0x6e, 0x65, 0x5f, 0x6e, 0x75, 0x6d, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x07, 0x6c, 0x69, 0x6e, 0x65, 0x4e, 0x75, 0x6d, 0x12, 0x1a, 0x0a,
	0x08, 0x65, 0x78, 0x70, 0x6f, 0x72, 0x74, 0x65, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x08, 0x52,
	0x08, 0x65, 0x78, 0x70, 0x6f, 0x72, 0x74, 0x65, 0x64, 0x42, 0x0c, 0x5a, 0x0a, 0x2e, 0x2f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_db_proto_rawDescOnce sync.Once
	file_db_proto_rawDescData = file_db_proto_rawDesc
)

func file_db_proto_rawDescGZIP() []byte {
	file_db_proto_rawDescOnce.Do(func() {
		file_db_proto_rawDescData = protoimpl.X.CompressGZIP(file_db_proto_rawDescData)
	})
	return file_db_proto_rawDescData
}

var file_db_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_db_proto_goTypes = []interface{}{
	(*Package)(nil), // 0: Package
	(*GoFile)(nil),  // 1: GoFile
	(*GoDef)(nil),   // 2: GoDef
}
var file_db_proto_depIdxs = []int32{
	1, // 0: Package.go_files:type_name -> GoFile
	2, // 1: GoFile.defs:type_name -> GoDef
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_db_proto_init() }
func file_db_proto_init() {
	if File_db_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_db_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Package); i {
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
		file_db_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GoFile); i {
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
		file_db_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GoDef); i {
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
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_db_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_db_proto_goTypes,
		DependencyIndexes: file_db_proto_depIdxs,
		MessageInfos:      file_db_proto_msgTypes,
	}.Build()
	File_db_proto = out.File
	file_db_proto_rawDesc = nil
	file_db_proto_goTypes = nil
	file_db_proto_depIdxs = nil
}

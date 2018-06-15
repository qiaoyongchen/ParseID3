package ParseID3

import (
    "bytes"
    "io"
    "os"
    "errors"
    "github.com/axgle/mahonia"
)

const (
    ID3V1_TAG                              =    "TAG"
    ID3V1_TAG_START_AT_DESC                =    128
    ID3V1_TAG_OFFSET                       =    3
    ID3V1_HEADER_OFFSET                    =    128
    
    ID3V2_TAG                              =    "ID3"
    ID3V2_TAG_START_AT                     =    0
    ID3V2_TAG_OFFSET                       =    3
    ID3V2_VERSION_START_AT                 =    3
    ID3V2_VERSION_OFFSET                   =    1
    ID3V2_VERSION1_START_AT                =    4
    ID3V2_VERSION1_OFFSET                  =    1
    ID3V2_FLAG_START_AT                    =    5
    ID3V2_FLAG_OFFSET                      =    1
    ID3V2_SIZE_START_AT                    =    6
    ID3V2_SIZE_OFFSET                      =    4
    ID3V2_HEADER_START_AT                  =    0
    ID3V2_HEADER_OFFSET                    =    10
    ID3V2_BODY_START_AT                    =    10
    ID3V2_FRAME_ID_OFFSET                  =    4
    ID3V2_FRAME_SIZE_OFFSET                =    4
    ID3V2_FRAME_FLAGS_OFFSET               =    2
    
    COVER_PHOTO_ID                         =    "APIC"
    
    UTF8_FLAG                              =    0x00
    UTF16_FLAG                             =    0x01
)

/*
    &错误列表
*/
var (
    ERROR_FILE_CANT_OPEN             =    errors.New("CAN'T OPEN FILE")
    ERROR_ID3V1_NOT_FOUND            =    errors.New("ID3V1 INFO NOT FOUND")
    ERROR_ID3V2_NOT_FOUND            =    errors.New("ID3V2 INFO NOT FOUND")
    ERROR_NEITHER_NOR_FOUND          =    errors.New("NEITHER ID3V1 NOR ID3V2 INFO WAS FOUND")
    ERROR_UNAMED_FRAME_ID            =    errors.New("UNAMED FRAME ID")
    ERROR_USE_ANOTHER_FUNC           =    errors.New("USE ANOTHER FUNCTION TO GET COVER : GetCover")
    ERROR_FRAME_ID_NOT_FOUND         =    errors.New("FRAME ID NOT FOUND")
)

/*
    FrameID 对应
*/
var FrameIDMap map[string]string = map[string]string {
    "AENC"    :    "音频加密技术",
    "APIC"    :    "附加描述",
    "COMM"    :    "注释",
    "COMR"    :    "广告",
    "ENCR"    :    "加密方法注册",
    "ETC0"    :    "事件时间编码",
    "GEOB"    :    "常规压缩对象",
    "GRID"    :    "组识别注册",
    "IPLS"    :    "复杂类别列表",
    "MCDI"    :    "音乐CD标识符",
    "MLLT"    :    "MPEG位置查找表格",
    "OWNE"    :    "所有权",
    "PRIV"    :    "私有",
    "PCNT"    :    "播放计数",
    "POPM"    :    "普通仪表",
    "POSS"    :    "位置同步",
    "RBUF"    :    "推荐缓冲区大小",
    "RVAD"    :    "音量调节器",
    "RVRB"    :    "混响",
    "SYLT"    :    "同步歌词或文本",
    "SYTC"    :    "同步节拍编码",
    "TALB"    :    "专辑",
    "TBPM"    :    "每分钟节拍数",
    "TCOM"    :    "作曲家",
    "TCON"    :    "流派",
    "TCOP"    :    "版权",
    "TDAT"    :    "日期",
    "TDLY"    :    "播放列表返录",
    "TENC"    :    "编码",
    "TEXT"    :    "歌词作者",
    "TFLT"    :    "文件类型",
    "TIME"    :    "时间",
    "TIT1"    :    "内容组描述",
    "TIT2"    :    "标题",
    "TIT3"    :    "副标题",
    "TKEY"    :    "最初关键字",
    "TLAN"    :    "语言",
    "TLEN"    :    "长度",
    "TMED"    :    "媒体类型",
    "TOAL"    :    "原唱片集",
    "TOFN"    :    "原文件名",
    "TOLY"    :    "原歌词作者",
    "TOPE"    :    "原艺术家",
    "TORY"    :    "最初发行年份",
    "TOWM"    :    "文件所有者",
    "TPE1"    :    "艺术家",
    "TPE2"    :    "乐队",
    "TPE3"    :    "指挥者",
    "TPE4"    :    "翻译",
    "TPOS"    :    "作品集部分",
    "TPUB"    :    "发行人",
    "TRCK"    :    "音轨(曲号)",
    "TRDA"    :    "录制日期",
    "TRSN"    :    "Intenet电台名称",
    "TRSO"    :    "Intenet电台所有者",
    "TSIZ"    :    "大小",
    "TSRC"    :    "ISRC(国际的标准记录代码)",
    "TSSE"    :    "编码使用的软件(硬件设置)",
    "TYER"    :    "年代",
    "TXXX"    :    "年度",
    "UFID"    :    "唯一的文件标识符",
    "USER"    :    "使用条款",
    "USLT"    :    "歌词",
    "WCOM"    :    "广告信息",
    "WCOP"    :    "版权信息",
    "WOAF"    :    "官方音频文件网页",
    "WOAR"    :    "官方艺术家网页",
    "WOAS"    :    "官方音频原始资料网页",
    "WORS"    :    "官方互联网无线配置首页",
    "WPAY"    :    "付款",
    "WPUB"    :    "出版商官方网页",
    "WXXX"    :    "用户定义的URL链接"}

/*
    ID3v1结构
*/
type ID3v1 struct {
    TAG            []byte                //"TAG"               [0,        2  ]
    Title          []byte                //                    [3,        32 ]
    Artist         []byte                //&                   [33,       62 ]
    Album          []byte                //&                   [63,       92 ]
    Year           []byte                //&                   [93,       96 ]
    Comments       []byte                //&                   [97,       124]
    Reserved       []byte                //&                   [125,      125]
    Track          []byte                //&                   [126,      126]
    Genre          []byte                //&                   [127,      127]
}

/*
    &标头信息占固定的十个字节
*/
type ID3v2Header struct {
    Tag            []byte              //"ID3"                 [0,  2]
    Ver            byte                //version               [3,  3]
    Ver1           byte                //version1              [4,  4]
    Flag           byte                //flag                  [5,  5]
    Size           []byte              //size                  [6,  9]
}

/*
    &获取标头长度
*/
func (this *ID3v2Header) GetSize () (int64) {
    return  int64(this.Size[0] & 0x7f)    *    0x200000    +
            int64(this.Size[1] & 0x7f)    *    0x4000      +
            int64(this.Size[2] & 0x7f)    *    0x80        +
            int64(this.Size[3] & 0x7f)
}

/*
    &标签体信息
*/
type ID3v2Frame struct {
    ID            []byte                //4 bytes
    Size          []byte                //4 bytes
    Flags         []byte                //2 bytes
    Content       []byte                //
}

/*
    seekBytes 获取文件字节数组
    f : 一个打开的文件句柄(有读权限)
    s : 开始位置
    o : 从开始位置起的后o个字节
*/
func fseek(f *os.File, s int64, o int64) []byte {
    var bb bytes.Buffer
    f.Seek(s, os.SEEK_SET)
    io.CopyN(&bb, f, o)
    return bb.Bytes()
}

/*
    &返回文件总长度
*/
func flen(f *os.File) int64 {
    l, _ := f.Seek(0, os.SEEK_END)
    return l
}

/*
    &解析ID3v2标签头信息
*/
func ParseID3v2Header(f *os.File) (*ID3v2Header, error) {
    id3h := &ID3v2Header{}
    tag := fseek(f, ID3V2_TAG_START_AT, ID3V2_TAG_OFFSET)
    if ID3V2_TAG != string(tag) {
        return nil, ERROR_ID3V2_NOT_FOUND
    }
    id3h.Tag = tag
    v := fseek(f, ID3V2_VERSION_START_AT, ID3V2_VERSION_OFFSET)
    id3h.Ver = v[0]
    v1 := fseek(f, ID3V2_VERSION1_START_AT, ID3V2_VERSION1_OFFSET)
    id3h.Ver1 = v1[0]
    fl := fseek(f, ID3V2_FLAG_START_AT, ID3V2_FLAG_OFFSET)
    id3h.Flag = fl[0]
    s := fseek(f, ID3V2_SIZE_START_AT, ID3V2_SIZE_OFFSET)
    id3h.Size = s
    return id3h, nil
}

/*
    &解析标签帧信息
    &返回标签帧结构,该标签帧末尾处位置,错误
*/
func ParseID3v2Frame(f *os.File, s int64) (*ID3v2Frame, int64, error) {
    id3f := &ID3v2Frame{}
    
    id3f.ID = fseek(f, s, ID3V2_FRAME_ID_OFFSET)
    s = s + ID3V2_FRAME_ID_OFFSET
    
    id3f.Size = fseek(f, s, ID3V2_FRAME_SIZE_OFFSET)
    s = s + ID3V2_FRAME_SIZE_OFFSET
    
    id3f.Flags = fseek(f, s, ID3V2_FRAME_FLAGS_OFFSET)
    s = s + ID3V2_FRAME_FLAGS_OFFSET
    
    id3f.Content = fseek(f, s, id3f.GetSize())
    s = s + id3f.GetSize()
    
    var exist bool
    if _, exist = FrameIDMap[string(id3f.ID)]; !exist {
        return nil, s, ERROR_UNAMED_FRAME_ID
    }
    return id3f, s, nil
}

/*
    &获取帧内容大小
*/
func (this *ID3v2Frame) GetSize() int64 {
    return  int64(this.Size[0])      *    0x1000000      +
            int64(this.Size[1])      *    0x10000        +
            int64(this.Size[2])      *    0x100          +
            int64(this.Size[3])
}

type ID3v2 struct {
    Header *ID3v2Header
    Frames map[string]*ID3v2Frame
}

/*
    &解析ID3v2信息
*/
func ParseID3v2(f *os.File) (*ID3v2, error) {
    id3v2 := &ID3v2{}
    id3v2.Frames = make(map[string]*ID3v2Frame)
    h, herror := ParseID3v2Header(f)
    if nil != herror {
        return nil, herror
    }
    id3v2.Header = h
    var s int64 = ID3V2_BODY_START_AT
    for {
        fr, l, frerror := ParseID3v2Frame(f, s)
        if nil != frerror {
            return id3v2, nil
        }
        id3v2.Frames[string(fr.ID)] = fr
        s = l
    }
    return id3v2, nil
}

func (this *ID3v2) GetCover() ([]byte, []byte) {
    if _, frexist := this.Frames[COVER_PHOTO_ID]; !frexist {
        return nil, nil
    }
    var i int64
    bytes := this.Frames[COVER_PHOTO_ID].Content
    
    //get mime type
    typebytes := make([]byte, 0)
    var firstbyte byte = UTF8_FLAG
    for i = 0; i < int64(len(bytes)); i = i + 1 {
        if i == 0 {
            if bytes[i] == UTF8_FLAG {
                continue
            } else if bytes[i] == UTF16_FLAG {
                firstbyte = UTF16_FLAG
            } else {
                typebytes = append(typebytes, bytes[i])
            }
        }
        if (0x06 == bytes[i]) && (0x00 == bytes[i + 1]) {
            break;
        }
        typebytes = append(typebytes, bytes[i])
    }
    if UTF16_FLAG == firstbyte {
        typebytes = encodeTranslate(typebytes, "utf16", "utf8")
    }
    
    binaryBytes := make([]byte, 0)
    inBinaryFlag := false
    imgType := ""
    for i = i; i < int64(len(bytes)); i = i + 1 {
        if false == inBinaryFlag {
            if 0xff == bytes[i + 0] && 0xd8 == bytes[i + 1] {
                binaryBytes = append(binaryBytes, bytes[i])
                inBinaryFlag = true
                imgType = "JPEG"
                continue
            }
            if 0x89 == bytes[i + 0] && 0x50 == bytes[i + 1] && 0x4e == bytes[i + 2] && 0x47 == bytes[i + 3] {
                binaryBytes = append(binaryBytes, bytes[i])
                inBinaryFlag = true
                imgType = "PNG"
                continue
            }
        }
        if true == inBinaryFlag && "" != imgType {
            binaryBytes = append(binaryBytes, bytes[i])
            if imgType == "JPEG" && 0xff == bytes[i + 0] && 0xd9 == bytes[i + 1] {
                binaryBytes = append(binaryBytes, bytes[i + 1])
                break
            }
            if imgType == "PNG" && 0x4e == bytes[i + 0] && 0x44 == bytes[i + 1] && 0xae == bytes[i + 2] && 
                0x42 == bytes[i + 3] && 0x60 == bytes[i + 4] && 0x82 == bytes[i + 5] {
                binaryBytes = append(binaryBytes, bytes[i + 1])
                binaryBytes = append(binaryBytes, bytes[i + 2])
                binaryBytes = append(binaryBytes, bytes[i + 3])
                binaryBytes = append(binaryBytes, bytes[i + 4])
                binaryBytes = append(binaryBytes, bytes[i + 5])
                break
            }
        }
    }
    return binaryBytes, typebytes
}

/*
    &获取帧体内容:string(utf8)
*/
func (this *ID3v2) GetFrameContent(fID string) ([]byte, error) {
    if _, ifinmap := FrameIDMap[fID]; !ifinmap {
        return nil, ERROR_UNAMED_FRAME_ID
    }
    if fID == COVER_PHOTO_ID {
        return nil, ERROR_USE_ANOTHER_FUNC
    }
    var fr *ID3v2Frame
    var ifinframes bool
    if fr, ifinframes = this.Frames[fID]; !ifinframes {
        return nil, ERROR_FRAME_ID_NOT_FOUND
    }
    contentbytes := fr.Content
    if contentbytes[0] == UTF8_FLAG {
        return contentbytes[1:], nil
    } else if contentbytes[0] == UTF16_FLAG {
        return encodeTranslate(contentbytes[1:], "utf16", "utf8"), nil
    } else {
        return contentbytes, nil
    }
}

/*
    &编码转换
*/
func encodeTranslate(source []byte, sourceEncode string, targetEncode string) []byte {
    s := mahonia.NewDecoder(sourceEncode).ConvertString(string(source))
    _, t, _ := mahonia.NewDecoder(targetEncode).Translate([]byte(s), true)
    return t
}


package mc_util

import (
	json "encoding/json"
	time "time"
	net "net"
	fmt "fmt"
)

type ServerInfo struct{
	Description string
	Max_player uint32
	Online_player uint32
	Players []map[string]string
	Version string
	Favicon string
	Delay int64
}

func Ping(host string, port uint16)(info *ServerInfo, err error){
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	r := NewReader(conn)
	w := NewWriter(conn)

	starttime := time.Now()

	w.BeginWritePacket(0x00)
	w.WriteVarInt32(0)
	w.WriteString(host)
	w.WriteUint16(port)
	w.WriteVarInt32(1)
	_, _, err = w.EndWritePacket()
	if err != nil {
		return nil, err
	}

	w.BeginWritePacket(0x00)
	_, _, err = w.EndWritePacket()
	if err != nil {
		return nil, err
	}

	_, err = r.BeginReadPacket()
	if err != nil {
		return nil, err
	}

	var jstr string
	jstr, _, err = r.ReadString()
	if err != nil {
		return nil, err
	}
	r.EndReadPacket()

	endtime := time.Now()
	delay := (int64)(endtime.Sub(starttime) / 1000000)

	obj := make(map[string]interface{})
	json.Unmarshal(([]byte)(jstr), &obj)
	players := obj["players"].(map[string]interface{})
	version := obj["version"].(map[string]interface{})
	return &ServerInfo{
		Description: parseDescStr(obj["description"].(map[string]interface{})),
		Max_player: (uint32)(players["max"].(float64)),
		Online_player: (uint32)(players["online"].(float64)),
		Players: players["sample"].([]map[string]string),
		Version: version["name"].(string),
		Favicon: obj["favicon"].(string),
		Delay: delay,
	}, nil
}

func parseDesc(buf *bytes.Buffer, descmap map[string]interface{}){
	buf.WriteString(descmap["text"].(string))
	extra0, eok := descmap["extra"]
	if !eok {
		return
	}
	extras := util.JsonToArrMap(extra0)
	for _, e := range extras {
		parseDesc(buf, e)
	}
}

func parseDescStr(descmap map[string]interface{})(string){
	buf := bytes.NewBuffer([]byte{})
	parseDesc(buf, descmap)
	return buf.String()
}



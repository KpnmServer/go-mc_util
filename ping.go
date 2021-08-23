
package mc_util

import (
	bytes "bytes"
	json "encoding/json"
	time "time"
	net "net"
	fmt "fmt"
)

type ServerStatus struct{
	Description string
	Max_player uint32
	Online_player uint32
	Players []map[string]interface{}
	Version string
	Favicon string
	Delay int64
}

func Ping(host string, port uint16)(status *ServerStatus, err error){
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", host, port), time.Second * 5)
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
	player_info := obj["players"].(map[string]interface{})
	version := obj["version"].(map[string]interface{})
	var players = []map[string]interface{}{}
	if player_info["sample"] != nil {
		players = make([]map[string]interface{}, 0, len(player_info["sample"].([]interface{})))
		for _, p := range player_info["sample"].([]interface{}){
			players = append(players, p.(map[string]interface{}))
		}
	}
	return &ServerStatus{
		Description: parseDescStr(obj["description"].(map[string]interface{})),
		Max_player: (uint32)(player_info["max"].(float64)),
		Online_player: (uint32)(player_info["online"].(float64)),
		Players: players,
		Version: version["name"].(string),
		Favicon: obj["favicon"].(string),
		Delay: delay,
	}, nil
}

func parseDesc(buf *bytes.Buffer, descmap map[string]interface{}){
	buf.WriteString(descmap["text"].(string))
	extra, eok := descmap["extra"]
	if !eok {
		return
	}
	for _, e := range extra.([]interface{}) {
		parseDesc(buf, e.(map[string]interface{}))
	}
}

func parseDescStr(descmap map[string]interface{})(string){
	buf := bytes.NewBuffer([]byte{})
	parseDesc(buf, descmap)
	return buf.String()
}



package bencode

import (
	"fmt"
	"reflect"
	"testing"
)

func TestDecodeInt(t *testing.T) {
	type test struct {
		Test uint64 `bencode:"test"`
	}
	have := &test{}
	want := &test{1234567890}
	input := []byte("d4:testi1234567890ee")
	err := Decode(input, have)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(have, want) {
		t.Errorf("Struct not properly hidrated: wanted %v but have %v", want, have)
	}
}

func TestDecodeString(t *testing.T) {
	type test struct {
		Test string `bencode:"test"`
	}
	have := &test{}
	want := &test{"test"}
	input := []byte("d4:test4:teste")
	err := Decode(input, have)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(have, want) {
		t.Errorf("Struct not properly hidrated: wanted %v but have %v", want, have)
	}
}

func TestDecodeStringInt(t *testing.T) {
	type test struct {
		Foo string `bencode:"foo"`
		Bar uint64 `bencode:"bar"`
	}
	have := &test{}
	want := &test{"test", 1234567890}
	input := []byte("d3:bari1234567890e3:foo4:teste")
	err := Decode(input, have)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(have, want) {
		t.Errorf("Struct not properly hidrated: wanted %v but have %v", want, have)
	}
}

func TestDecodeMockTorrentFile(t *testing.T) {
	type info struct {
		Length uint64 `bencode:"length"`
		Name string `bencode:"name"`
		PieceLength uint64 `bencode:"piece length"`
	}
	type test struct {
		Announce string `bencode:"announce"`
		Comment string `bencode:"comment"`
		AnnounceList [][]string `bencode:"announce-list"`
		Info info `bencode:"info"`
	}
	have := &test{}
	want := &test{"https://torrent.ubuntu.com/announce", "Ubuntu CD releases.ubuntu.com", [][]string{{"https://torrent.ubuntu.com/announce"}, {"https://ipv6.torrent.ubuntu.com/announce"}}, info{4932407296, "ubuntu-23.04-desktop-amd64.iso", 262144}}
	input := []byte("d8:announce35:https://torrent.ubuntu.com/announce13:announce-listll35:https://torrent.ubuntu.com/announceel40:https://ipv6.torrent.ubuntu.com/announceee7:comment29:Ubuntu CD releases.ubuntu.com10:created by13:mktorrent 1.113:creation datei1681992794e4:infod6:lengthi4932407296e4:name30:ubuntu-23.04-desktop-amd64.iso12:piece lengthi262144eee")
	err := Decode(input, have)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(have, want) {
		fmt.Println(have.Info)
		fmt.Println(want.Info)
		t.Errorf("Struct not properly hidrated: wanted %+v but have %+v", want, have)
	}
}

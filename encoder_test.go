package bencode

import (
	"reflect"
	"testing"
)

func TestEcodeInt(t *testing.T) {
	type test struct {
		Test uint64 `bencode:"test"`
	}
	have, err := Encode(&test{1234567890})
	want := []byte("d4:testi1234567890ee")
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(have, want) {
		t.Errorf("Struct not properly encoded: wanted %s but have %s", want, have)
	}
}

func TestEncodeString(t *testing.T) {
	type test struct {
		Test string `bencode:"test"`
	}
	have, err := Encode(&test{"test"})
	want := []byte("d4:test4:teste")
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(have, want) {
		t.Errorf("Struct not properly encoded: wanted %s but have %s", want, have)
	}
}

func TestEncodeStringInt(t *testing.T) {
	type test struct {
		TestString string `bencode:"teststring"`
		TestInt uint64 `bencode:"testint"`
	}
	have, err := Encode(&test{"test", 1234567890})
	want := []byte("d10:teststring4:test7:testinti1234567890ee")
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(have, want) {
		t.Errorf("Struct not properly encoded: wanted %s but have %s", want, have)
	}
}

func TestEncodeMockTorrentFile(t *testing.T) {
	type info struct {
		Length uint64 `bencode:"length"`
		Name string `bencode:"name"`
		PieceLength uint64 `bencode:"piece length"`
	}
	type test struct {
		Announce string `bencode:"announce"`
		AnnounceList [][]string `bencode:"announce-list"`
		Comment string `bencode:"comment"`
		CreatedBy string `bencode:"created by"`
		CreationDate uint64 `bencode:"creation date"`
		Info info `bencode:"info"`
	}
	have, err := Encode(&test{"https://torrent.ubuntu.com/announce", [][]string{{"https://torrent.ubuntu.com/announce"}, {"https://ipv6.torrent.ubuntu.com/announce"}}, "Ubuntu CD releases.ubuntu.com", "mktorrent 1.1", 1681992794, info{4932407296, "ubuntu-23.04-desktop-amd64.iso", 262144}})
	want := []byte("d8:announce35:https://torrent.ubuntu.com/announce13:announce-listll35:https://torrent.ubuntu.com/announceel40:https://ipv6.torrent.ubuntu.com/announceee7:comment29:Ubuntu CD releases.ubuntu.com10:created by13:mktorrent 1.113:creation datei1681992794e4:infod6:lengthi4932407296e4:name30:ubuntu-23.04-desktop-amd64.iso12:piece lengthi262144eee")
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(have, want) {
		t.Errorf("Struct not properly encoded: wanted %s but have %s", want, have)
	}
}

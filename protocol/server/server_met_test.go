package server

import (
	"bytes"
	"encoding/binary"
	"os"
	"path/filepath"
	"testing"
)

func TestLoadServerMetFixtures(t *testing.T) {
	fixtures := []string{
		filepath.Join("..", "..", "..", "jed2k", "core", "src", "main", "resources", "server.met"),
		filepath.Join("..", "..", "..", "jed2k", "core", "src", "main", "resources", "server2.met"),
	}
	for _, path := range fixtures {
		if _, err := os.Stat(path); err != nil {
			if os.IsNotExist(err) {
				t.Skipf("fixture %s not available", path)
			}
			t.Fatalf("stat %s: %v", path, err)
		}
		met, err := LoadServerMet(path)
		if err != nil {
			t.Fatalf("load %s: %v", path, err)
		}
		if len(met.Servers) == 0 {
			t.Fatalf("expected entries in %s", path)
		}
		if len(met.Addresses()) == 0 {
			t.Fatalf("expected addresses in %s", path)
		}
	}
}

func TestServerMetRoundTripAndGetters(t *testing.T) {
	met := NewServerMet()

	ipEntry, err := NewServerMetEntryFromIP("192.168.0.9", 5600, "Test server name", "Test descr")
	if err != nil {
		t.Fatalf("ip entry: %v", err)
	}
	hostEntry, err := NewServerMetEntryFromHost("mule.org", 45567, "Name2", "")
	if err != nil {
		t.Fatalf("host entry: %v", err)
	}
	met.AddServer(ipEntry)
	met.AddServer(hostEntry)

	var payload bytes.Buffer
	if err := met.Put(&payload); err != nil {
		t.Fatalf("put: %v", err)
	}

	var parsed ServerMet
	if err := parsed.Get(bytes.NewReader(payload.Bytes())); err != nil {
		t.Fatalf("get: %v", err)
	}
	if got := parsed.Servers[0].Name(); got != "Test server name" {
		t.Fatalf("unexpected first name: %q", got)
	}
	if got := parsed.Servers[0].Description(); got != "Test descr" {
		t.Fatalf("unexpected first description: %q", got)
	}
	if got := parsed.Servers[0].Host(); got != "192.168.0.9" {
		t.Fatalf("unexpected first host: %q", got)
	}
	if got := parsed.Servers[1].Name(); got != "Name2" {
		t.Fatalf("unexpected second name: %q", got)
	}
	if got := parsed.Servers[1].Description(); got != "" {
		t.Fatalf("unexpected second description: %q", got)
	}
	if got := parsed.Servers[1].Host(); got != "mule.org" {
		t.Fatalf("unexpected second host: %q", got)
	}
	if got := parsed.Addresses(); len(got) != 2 || got[0] != "192.168.0.9:5600" || got[1] != "mule.org:45567" {
		t.Fatalf("unexpected addresses: %#v", got)
	}
}

func TestParseServerMetRejectsInvalidHeader(t *testing.T) {
	payload := []byte{'<', 'h', 't', 'm', 'l'}

	if _, err := ParseServerMet(payload); err == nil {
		t.Fatal("expected invalid header error")
	}
}

func TestParseServerMetRejectsImpossibleEntryCount(t *testing.T) {
	var payload bytes.Buffer
	if err := payload.WriteByte(serverMetHeader); err != nil {
		t.Fatalf("write header: %v", err)
	}
	if err := binary.Write(&payload, binary.LittleEndian, uint32(1)); err != nil {
		t.Fatalf("write count: %v", err)
	}

	if _, err := ParseServerMet(payload.Bytes()); err == nil {
		t.Fatal("expected invalid entry count error")
	}
}

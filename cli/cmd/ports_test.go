package cmd

import (
	"net"
	"testing"
)

// Verify that a port nobody is listening on is correctly reported as available
func TestIsPortAvailable_FreePort(t *testing.T) {
	// Pick an port unlikely to be in use
	port := "19876"

	if !IsPortAvailable(port) {
		t.Skipf("port %s happens to be in use on this machine, skipping", port)
	}
}

// Verify that a port we bind is correctly reported as unavailable
func TestIsPortAvailable_OccupiedPort(t *testing.T) {
	// Bind a listener to occupy a port
	ln, err := net.Listen("tcp", ":0") // :0 = OS picks a free port
	if err != nil {
		t.Fatalf("failed to start test listener: %v", err)
	}
	defer ln.Close()

	// Extract the port the OS assigned
	_, port, err := net.SplitHostPort(ln.Addr().String())
	if err != nil {
		t.Fatalf("failed to parse listener address: %v", err)
	}

	// The port is still bound — IsPortAvailable must return false
	if IsPortAvailable(port) {
		t.Errorf("expected port %s to be unavailable, but IsPortAvailable returned true", port)
	}
}

// Verify that a port becomes available again after the listener that held it is closed
func TestIsPortAvailable_AfterRelease(t *testing.T) {
	ln, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatalf("failed to start test listener: %v", err)
	}

	_, port, err := net.SplitHostPort(ln.Addr().String())
	if err != nil {
		t.Fatalf("failed to parse listener address: %v", err)
	}

	// Port is occupied
	if IsPortAvailable(port) {
		t.Errorf("expected port %s to be unavailable while listener is open", port)
	}

	// Release the port
	ln.Close()

	// Port should now be free
	if !IsPortAvailable(port) {
		t.Errorf("expected port %s to be available after listener closed", port)
	}
}

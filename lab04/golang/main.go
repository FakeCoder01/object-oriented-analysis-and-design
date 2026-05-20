package main

import (
	"fmt"
	"net"
	"strings"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type Record struct {
	Key   string
	Value string
}

type Gateway struct {
	Host string
	Port int
	Conn net.Conn
}

func (g *Gateway) Connect() error {
	address := fmt.Sprintf("%s:%d", g.Host, g.Port)
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return err
	}
	g.Conn = conn
	return nil
}

func (g *Gateway) Disconnect() {
	if g.Conn != nil {
		g.Conn.Close()
	}
}

func (g *Gateway) ExecuteCommand(cmd string) (string, error) {
	if g.Conn == nil {
		return "", fmt.Errorf("not connected")
	}

	// send cmd
	_, err := fmt.Fprintf(g.Conn, cmd)
	if err != nil {
		return "", err
	}

	// read resp (w/o '\n')
	buf := make([]byte, 1024)
	n, err := g.Conn.Read(buf)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(buf[:n])), nil
}

type IdentityMap struct {
	cache   map[string]*Record
	gateway *Gateway
	mu      sync.Mutex
	logFunc func(string)
}

func NewIdentityMap(gw *Gateway, logger func(string)) *IdentityMap {
	return &IdentityMap{
		cache:   make(map[string]*Record),
		gateway: gw,
		logFunc: logger,
	}
}

func (im *IdentityMap) Find(key string) *Record {
	im.mu.Lock()
	defer im.mu.Unlock()

	if record, exists := im.cache[key]; exists {
		im.logFunc(fmt.Sprintf("[IDENTITY MAP] CACHE HIT! Returning existing pointer %p for key: %s", record, key))
		return record
	}

	im.logFunc("[IDENTITY MAP] CACHE MISS. Fetching from MeowDB server...")
	response, err := im.gateway.ExecuteCommand("select " + key)
	if err == nil && !strings.Contains(response, "not found") {
		record := &Record{Key: key, Value: response}
		im.cache[key] = record
		im.logFunc(fmt.Sprintf("[IDENTITY MAP] Object instantiated at %p and mapped.", record))
		return record
	}
	return nil
}


func (im *IdentityMap) Save(command, key, value string) {
	im.mu.Lock()
	defer im.mu.Unlock()

	response, err := im.gateway.ExecuteCommand(fmt.Sprintf("%s %s %s", command, key, value))
	if err == nil && !strings.Contains(response, "not found") {
		if record, exists := im.cache[key]; exists {
			record.Value = value
			im.logFunc(fmt.Sprintf("[IDENTITY MAP] Existing object %p updated.", record))
		} else {
			record := &Record{Key: key, Value: value}
			im.cache[key] = record
			im.logFunc(fmt.Sprintf("[IDENTITY MAP] New object %p added.", record))
		}
	}
}

func (im *IdentityMap) Delete(key string) {
	im.mu.Lock()
	defer im.mu.Unlock()

	response, err := im.gateway.ExecuteCommand("delete " + key)
	if err == nil && !strings.Contains(response, "not found") {
		delete(im.cache, key)
		im.logFunc("[IDENTITY MAP] Object removed from map.")
	} else {
		im.logFunc("[DATABASE] Key not found to delete.")
	}
}

func (im *IdentityMap) GetAllCached() []*Record {
	im.mu.Lock()
	defer im.mu.Unlock()
	var records []*Record
	for _, rec := range im.cache {
		records = append(records, rec)
	}
	return records
}


func main() {
	a := app.New()
	w := a.NewWindow("MeowDB - Go Identity Map")
	w.Resize(fyne.NewSize(800, 500))

	var gateway *Gateway
	var identityMap *IdentityMap
	connected := false

	hostEntry := widget.NewEntry()
	hostEntry.SetText("meowdb")

	keyEntry := widget.NewEntry()
	valueEntry := widget.NewEntry()

	console := widget.NewMultiLineEntry()
	console.Disable() // read-only

	logFunc := func(msg string) {
		console.SetText(console.Text + msg + "\n")
		console.CursorRow = len(strings.Split(console.Text, "\n"))
	}

	var cachedData []*Record
	list := widget.NewTable(
		func() (int, int) { return len(cachedData), 3 },
		func() fyne.CanvasObject { return widget.NewLabel("Template Label Width") },
		func(i widget.TableCellID, o fyne.CanvasObject) {
			label := o.(*widget.Label)
			rec := cachedData[i.Row]
			switch i.Col {
			case 0: label.SetText(fmt.Sprintf("%p", rec))
			case 1: label.SetText(rec.Key)
			case 2: label.SetText(rec.Value)
			}
		},
	)
	list.SetColumnWidth(0, 150)
	list.SetColumnWidth(1, 150)
	list.SetColumnWidth(2, 200)

	refreshTable := func() {
		if identityMap != nil {
			cachedData = identityMap.GetAllCached()
			list.Refresh()
		}
	}

	var btnConnect *widget.Button
	btnConnect = widget.NewButton("Connect via TCP", func() {
		if !connected {
			gateway = &Gateway{Host: hostEntry.Text, Port: 6969}
			err := gateway.Connect()
			if err != nil {
				logFunc("Connection failed: " + err.Error())
				return
			}
			identityMap = NewIdentityMap(gateway, logFunc)
			connected = true
			btnConnect.SetText("Disconnect")
			logFunc("Connected to MeowDB.")
		} else {
			gateway.Disconnect()
			connected = false
			btnConnect.SetText("Connect via TCP")
			logFunc("Disconnected.")
		}
	})

	btnSelect := widget.NewButton("Select", func() {
		if !connected { return }
		rec := identityMap.Find(keyEntry.Text)
		if rec != nil {
			valueEntry.SetText(rec.Value)
			refreshTable()
		} else {
			logFunc("[DATABASE] Key not found.")
		}
	})

	btnInsert := widget.NewButton("Insert / Update", func() {
		if !connected { return }
		identityMap.Save("insert", keyEntry.Text, valueEntry.Text)
		refreshTable()
	})

	btnDelete := widget.NewButton("Delete", func() {
		if !connected { return }
		identityMap.Delete(keyEntry.Text)
		keyEntry.SetText("")
		valueEntry.SetText("")
		refreshTable()
	})

	topPanel := container.NewHBox(widget.NewLabel("Host:"), hostEntry, btnConnect)
	leftControls := container.NewVBox(
		widget.NewLabel("Key:"), keyEntry,
		widget.NewLabel("Value:"), valueEntry,
		btnSelect, btnInsert, btnDelete,
	)

	mainLayout := container.NewBorder(
		topPanel,
		container.NewVScroll(console),
		leftControls,
		nil,
		list,
	)

	w.SetContent(mainLayout)
	w.ShowAndRun()
}

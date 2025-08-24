package portscanner

import (
	"context"
	"encoding/csv"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/alexeyco/simpletable"
	"golang.org/x/sync/semaphore"
)

type StatusCode string

const (
	StatusOpen   StatusCode = "Open"
	StatusClosed StatusCode = "Closed"
)

type service struct {
	port        int
	service     string
	description string
}

type ServiceMap struct {
	services map[int]service
}

func newServiceMap() *ServiceMap {
	return &ServiceMap{
		make(map[int]service),
	}
}

type portData struct {
	number       int
	statusCode   StatusCode
	service      string
	descpription string
}

type portScanner struct {
	ip    string
	lock  *semaphore.Weighted
	ports []portData
}

func NewPortScanner(ip string, weight *semaphore.Weighted) *portScanner {
	return &portScanner{
		ip:    ip,
		lock:  weight,
		ports: make([]portData, 0),
	}
}

func getProtocolFromPort(pathtofile string) (*ServiceMap, error) {

	servicemp := newServiceMap()

	file, err := os.Open(pathtofile)
	if err != nil {
		return &ServiceMap{}, err
	}
	defer file.Close()

	reader := csv.NewReader(file)

	records, err := reader.ReadAll()
	if err != nil {
		return &ServiceMap{}, err
	}

	for _, row := range records {

		port, err := strconv.Atoi(row[0])
		if err != nil {
			return &ServiceMap{}, err
		}

		servicemp.services[port] = service{
			port:        port,
			service:     row[1],
			description: row[2],
		}
	}

	return servicemp, nil
}

func scanPort(ip string, targetedPort int, timeout time.Duration, ser ServiceMap) portData {
	target := fmt.Sprintf("%s:%d", ip, targetedPort)
	conn, err := net.DialTimeout("tcp", target, timeout)

	if err != nil {
		if strings.Contains(err.Error(), "too many open files") {
			time.Sleep(timeout)
			scanPort(ip, targetedPort, timeout, ser)
		} else {
			return portData{
				number:     targetedPort,
				statusCode: StatusClosed,
			}
		}
	}

	portData := portData{}
	portData.number = targetedPort
	portData.statusCode = StatusOpen
	serv, ok := ser.services[targetedPort]
	if ok {
		portData.service = serv.service
		portData.descpription = serv.description
	} else {
		portData.service = ""
	}

	conn.Close()
	return portData
}

func (p *portScanner) Scan(limit int, timeout time.Duration) {
	wg := sync.WaitGroup{}
	defer wg.Wait()

	serviceMap, err := getProtocolFromPort("listofports.csv")
	if err != nil {
		return
	}

	for port := 0; port <= limit; port++ {
		wg.Add(1)
		p.lock.Acquire(context.TODO(), 1)

		go func(port int) {
			defer p.lock.Release(1)
			defer wg.Done()
			scannedPort := scanPort(p.ip, port, timeout, *serviceMap)

			p.ports = append(p.ports, scannedPort)
		}(port)
	}
}

func (p *portScanner) ShowOpenPorts() {
	table := simpletable.New()

	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "#"},
			{Align: simpletable.AlignCenter, Text: "PORT"},
			{Align: simpletable.AlignCenter, Text: "STATUS"},
			{Align: simpletable.AlignCenter, Text: "SERVICE"},
			{Align: simpletable.AlignCenter, Text: "DESCRIPTION"},
		},
	}
	count := 0
	for _, row := range p.ports {
		if row.statusCode == StatusOpen {
			count++
			r := []*simpletable.Cell{
				{Align: simpletable.AlignCenter, Text: fmt.Sprintf("%d", count)},
				{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%d", row.number)},
				{Align: simpletable.AlignCenter, Text: string(row.statusCode)},
				{Align: simpletable.AlignLeft, Text: row.service},
				{Align: simpletable.AlignLeft, Text: row.descpription},
			}

			table.Body.Cells = append(table.Body.Cells, r)

		}
	}

	table.SetStyle(simpletable.StyleCompactClassic)
	fmt.Println(table)

}

package comm

import (
    "fmt"
    "github.com/go-ole/go-ole"
    "github.com/go-ole/go-ole/oleutil"
    "log"
    "regexp"
    "strconv"
)

// NOTE:  This function will panic if called with high frequency (we should throttle the polling - probably caused by slow OLE calls that can't be parallel.. sigh ;))
func GetSerialPortList() ([]CommPort, error) {
    err := ole.CoInitializeEx(0, ole.COINIT_MULTITHREADED)
    if err != nil {
        log.Fatal("Init error: ", err)
    }
    defer ole.CoUninitialize()

    unknown, _ := oleutil.CreateObject("WbemScripting.SWbemLocator")
    defer unknown.Release()

    wmi, _ := unknown.QueryInterface(ole.IID_IDispatch)
    defer wmi.Release()

    serviceRaw, _ := oleutil.CallMethod(wmi, "ConnectServer")
    service := serviceRaw.ToIDispatch()
    defer service.Release()

    query := "SELECT * FROM Win32_PnPEntity WHERE ConfigManagerErrorCode = 0 and Name like '%(COM%'"
    queryResult, err2 := oleutil.CallMethod(service, "ExecQuery", query)

    commPorts := []CommPort{}

    if err2 != nil {
        log.Printf("Error from oleutil.CallMethod: ", err2)
        return nil, err2
    }

    result := queryResult.ToIDispatch()
    defer result.Release()

    countVar, _ := oleutil.GetProperty(result, "Count")
    count := int(countVar.Val)

    // Should we filter on VID/PID here or later?
    for i := 0; i < count; i++ {
        itemRaw, _ := oleutil.CallMethod(result, "ItemIndex", i)
        item := itemRaw.ToIDispatch()
        defer item.Release()

        displayName, _ := oleutil.GetProperty(item, "Name")

        fmt.Printf("Port found: %s\n", displayName.ToString())

        re := regexp.MustCompile("\\((COM[0-9]+)\\)").FindAllStringSubmatch(displayName.ToString(), 1)

        var path string = ""

        if re != nil && len(re[0]) > 1 {
            path = re[0][1]
        }

        fmt.Printf("Path: %v\n", path)

        deviceId, _ := oleutil.GetProperty(item, "DeviceID")

        re = regexp.MustCompile("ID_(....)").FindAllStringSubmatch(deviceId.ToString(), 2)

        var VID, PID uint16 = 0, 0

        if re != nil && len(re) == 2 {
            if len(re[0]) > 1 {
                val, _ := strconv.ParseUint(re[0][1], 16, 16)
                VID = uint16(val)
            }
            if len(re[1]) > 1 {
                val, _ := strconv.ParseUint(re[1][1], 16, 16)
                PID = uint16(val)
            }
            fmt.Printf("VID: %v, PID: %v\n", VID, PID)
        }

        commPorts = append(commPorts, CommPort{Path: path, VendorId: VID, ProductId: PID, DisplayName: displayName.ToString()})

    }

    return commPorts, err
}

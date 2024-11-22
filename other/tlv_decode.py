import base64
from datetime import datetime

def fmtTimestamp(timestamp):
    dt = datetime.fromtimestamp(timestamp)
    return dt.strftime("%Y-%m-%d %H:%M:%S")

# Little endian byte to integer
def bytesToIntLittleEndian(byteArray):
    val = 0;
    for i in range(len(byteArray)): 
        val = val | byteArray[i] << i*8
    return val


# 解析以4347开头的hex数据
def tlvUnpack(byteArray):
    cmd = byteArray[2:3].hex()
    length = bytesToIntLittleEndian(byteArray[3:5])
    payload = byteArray[5:5+length]
    productId = 0

    index = 0
    subPackList = [];
    while index < length:
        key = payload[index:index+1].hex()
        subLen = bytesToIntLittleEndian(payload[index+1:index+3])
        subPayload = payload[index+3:index+3+subLen]
        index = index+3+subLen
        subPack = {
            'key': key,
            'len': subLen,
            'payload': subPayload,
        }

        if key == '38':
            productId = subPayload[0]

        #print('key:',key,'payload:',subPayload.hex())
        subPackList.append(subPack)

    return {
        'cmd': cmd,
        'productId': productId,
        'length': length,
        'subPackList': subPackList,
    }

def decodeTHdata(byteArray):
    th = bytesToIntLittleEndian(byteArray[0:3])
    temperature = ((th >> 12) - 500) / 10
    humidity = (th & 0xFFF) / 10
    pressure = bytesToIntLittleEndian(byteArray[3:5])
    battery = byteArray[5]

    return {
        'dataType': 'data',
        'timestamp': 0,
        'time': '',
        'temperature': temperature,
        'humidity': humidity,
        'pressure': pressure,
        'battery': battery,
    }

def decodeRealTimeData(byteArray,productId = 0):
    timestamp = bytesToIntLittleEndian(byteArray[0:4])
    realtimeData = decodeTHdata(byteArray[4:])
    rssi = byteArray[10]
    if rssi >= 128:
        rssi = rssi - 256
    
    realtimeData['dataType'] = 'event'
    realtimeData['timestamp'] = timestamp
    realtimeData['time'] = fmtTimestamp(timestamp)

    realtimeData['rssi'] = rssi
    return realtimeData

def decodeHistoryData(byteArray,productId = 0):
    timestamp = bytesToIntLittleEndian(byteArray[0:4])
    duration = bytesToIntLittleEndian(byteArray[4:6])

    historyDataList = []
    packLen = 6
    index = 6
    i = 0
    while index < len(byteArray):
        historyPack = byteArray[index:index+packLen]
        historyData = decodeTHdata(historyPack)
        historyData['timestamp'] = timestamp + duration * i
        historyData['dataType'] = 'data'
        historyData['time'] = fmtTimestamp(historyData['timestamp'])

        historyDataList.append(historyData)
        i += 1
        index += packLen

    return historyDataList


def tlvDecode(byteArray):
    unpackData = tlvUnpack(byteArray)
    outData = {'productId': unpackData['productId']}
    
    for subPack in unpackData['subPackList']: 
        if subPack['key'] == '14':
            realtimeData = decodeRealTimeData(subPack['payload'], unpackData['productId'])
            outData['sensorData'] = [realtimeData];
        
        if subPack['key'] == '03':
            historyeData = decodeHistoryData(subPack['payload'], unpackData['productId'])
            outData['sensorData'] = historyeData;

        if subPack['key'] == '11':
            outData['version'] = subPack['payload'].decode("utf-8")
        
        if subPack['key'] == '34':
            outData['versionModel'] = subPack['payload'].decode("utf-8")

        if subPack['key'] == '35':
            outData['versionMcu'] = subPack['payload'].decode("utf-8")

        if subPack['key'] == '04':
            reportInterval = bytesToIntLittleEndian(subPack['payload'])
            outData['reportInterval'] = reportInterval
        if subPack['key'] == '05':
            collectInterval = bytesToIntLittleEndian(subPack['payload'])
            outData['collectInterval'] = collectInterval

    return outData

# test
if __name__ == '__main__':
    # decode hex data
    src = '43473442003802002900110500322e302e36220400303030302c01000067040004000000340500312e392e35350500322e302e361d010001140c00a82b0f6707332e00003ae6006109'
    bs = bytes.fromhex(src)
    out = tlvDecode(bs)
    print(out)

    # # decode bas64 data
    # src = 'Q0cxMAA4AgApAB0BAAEDJAAc4j9nhAOsQS4AADarMS4AADatMS4AADauMS4AADauMS4AADYYCg=='
    # bs = base64.b64decode(src)
    # out = tlvDecode(bs)
    # print(out)

    # # decode base64 setting
    # src = 'Q0cyCgAFAgBYAgQCADwAaQE='
    # bs = base64.b64decode(src)
    # out = tlvDecode(bs)
    # print(out)
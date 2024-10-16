package other; 

// FIXME rename package

import java.util.*;

public class TLVDecoder {

    // 类：SubPack，表示子包
    public static class SubPack {
        private String key;
        private int len;
        private byte[] payload;

        public SubPack(String key, int len, byte[] payload) {
            this.key = key;
            this.len = len;
            this.payload = payload;
        }

        public String getKey() {
            return key;
        }

        public int getLen() {
            return len;
        }

        public byte[] getPayload() {
            return payload;
        }

        @Override
        public String toString() {
            return "SubPack{" +
                    "key='" + key + '\'' +
                    ", len=" + len +
                    ", payload=" + bytesToHex(payload) +
                    '}';
        }
    }

    // 类：TlvUnpackResult，表示 TLV 解包结果
    public static class TlvUnpackResult {
        private String cmd;
        private int length;
        private List<SubPack> subPackList;

        public TlvUnpackResult(String cmd, int length, List<SubPack> subPackList) {
            this.cmd = cmd;
            this.length = length;
            this.subPackList = subPackList;
        }

        public String getCmd() {
            return cmd;
        }

        public int getLength() {
            return length;
        }

        public List<SubPack> getSubPackList() {
            return subPackList;
        }

        @Override
        public String toString() {
            return "TlvUnpackResult{" +
                    "cmd='" + cmd + '\'' +
                    ", length=" + length +
                    ", subPackList=" + subPackList +
                    '}';
        }
    }

    // 类：DecodedRealTimeData，表示解码后的实时数据
    public static class DecodedRealTimeData {
        private String dataType;
        private int timestamp;
        private double temperature;
        private double humidity;
        private int pressure;
        private int battery;
        private int rssi;

        public DecodedRealTimeData(String dataType, int timestamp, double temperature, double humidity, int pressure, int battery, int rssi) {
            this.dataType = dataType;
            this.timestamp = timestamp;
            this.temperature = temperature;
            this.humidity = humidity;
            this.pressure = pressure;
            this.battery = battery;
            this.rssi = rssi;
        }

        public String getDataType() {
            return dataType;
        }

        public int getTimestamp() {
            return timestamp;
        }

        public double getTemperature() {
            return temperature;
        }

        public double getHumidity() {
            return humidity;
        }

        public int getPressure() {
            return pressure;
        }

        public int getBattery() {
            return battery;
        }

        public int getRssi() {
            return rssi;
        }

        @Override
        public String toString() {
            return "DecodedRealTimeData{" +
                    "dataType='" + dataType + '\'' +
                    ", timestamp=" + timestamp +
                    ", temperature=" + temperature +
                    ", humidity=" + humidity +
                    ", pressure=" + pressure +
                    ", battery=" + battery +
                    ", rssi=" + rssi +
                    '}';
        }
    }

    // 方法：将十六进制字符串转换为字节数组
    public static byte[] hexStringToByteArray(String s) {
        int len = s.length();
        if (len % 2 != 0) {
            throw new IllegalArgumentException("十六进制字符串长度必须为偶数");
        }
        byte[] data = new byte[len / 2];
        for(int i = 0; i < len; i += 2){
            data[i/2] = (byte)((Character.digit(s.charAt(i), 16) << 4)
                                 + Character.digit(s.charAt(i+1), 16));
        }
        return data;
    }

    // 方法：将字节数组转换为有符号整数（小端序）
    public static int bytesToIntLittleEndian(byte[] byteArray) {
        int val = 0;
        for (int i = 0; i < byteArray.length; i++) {
            val |= (byteArray[i] & 0xFF) << (i * 8);
        }
        return val;
    }

    // 方法：将字节数组转换为十六进制字符串
    public static String bytesToHex(byte[] bytes) {
        StringBuilder sb = new StringBuilder();
        for(byte b : bytes){
            sb.append(String.format("%02x", b));
        }
        return sb.toString();
    }

    // 方法：解包 TLV 数据
    public static TlvUnpackResult tlvUnpack(byte[] byteArray) {
        if (byteArray.length < 5) {
            throw new IllegalArgumentException("字节数组长度不足以解包 TLV 数据");
        }

        String cmd = String.format("%02x", byteArray[2]); // byteArray[2:3].hex()
        int length = bytesToIntLittleEndian(Arrays.copyOfRange(byteArray, 3, 5)); // byteArray[3:5]

        if (byteArray.length < 5 + length) {
            throw new IllegalArgumentException("字节数组长度不足以提取 payload");
        }

        byte[] payload = Arrays.copyOfRange(byteArray, 5, 5 + length);

        int index = 0;
        List<SubPack> subPackList = new ArrayList<>();
        while (index < length) {
            if (index + 1 > payload.length) {
                throw new IllegalArgumentException("子包格式错误：无法提取 key");
            }
            String key = String.format("%02x", payload[index]); // payload[index:index+1].hex()
            index += 1;

            if (index + 2 > payload.length) {
                throw new IllegalArgumentException("子包格式错误：无法提取 subLen");
            }
            int subLen = bytesToIntLittleEndian(Arrays.copyOfRange(payload, index, index + 2)); // payload[index:index+2]
            index += 2;

            if (index + subLen > payload.length) {
                throw new IllegalArgumentException("子包格式错误：subPayload 超出范围");
            }
            byte[] subPayload = Arrays.copyOfRange(payload, index, index + subLen);
            index += subLen;

            SubPack subPack = new SubPack(key, subLen, subPayload);
            subPackList.add(subPack);
        }

        return new TlvUnpackResult(cmd, length, subPackList);
    }

    // 方法：解码实时数据
    public static DecodedRealTimeData decodeRealTimeData(byte[] byteArray) {
        if (byteArray.length < 11) {
            throw new IllegalArgumentException("实时数据字节数组长度不足");
        }

        int timestamp = bytesToIntLittleEndian(Arrays.copyOfRange(byteArray, 0, 4));
        int th = bytesToIntLittleEndian(Arrays.copyOfRange(byteArray, 4, 7));
        double temperature = ((th >> 12) - 500) / 10.0;
        double humidity = (th & 0xFFF) / 10.0;
        int pressure = bytesToIntLittleEndian(Arrays.copyOfRange(byteArray, 7, 9));
        int battery = byteArray[9] & 0xFF; // 无符号
        int rssi = byteArray[10] & 0xFF; // 无符号
        if (rssi >= 128) {
            rssi -= 256;
        }

        return new DecodedRealTimeData("event", timestamp, temperature, humidity, pressure, battery, rssi);
    }

    // 方法：解码历史数据（暂未实现）
    public static void decodeHistoryData(byte[] byteArray) {
        // TODO: 实现历史数据的解码逻辑
    }

    // 方法：解析 TLV 数据并处理
    public static void tlvDecode(byte[] byteArray) {
        TlvUnpackResult unpackData = tlvUnpack(byteArray);

        for (SubPack subPack : unpackData.getSubPackList()) {
            switch (subPack.key) {
                case "14":
                    DecodedRealTimeData realtimeData = decodeRealTimeData(subPack.getPayload());
                    System.out.println(realtimeData);
                    break;
                
                case "03":
                    decodeHistoryData(subPack.getPayload());
                    break;
                
                default:
                    break;
            }
        }
    }

    // 主方法
    public static void main(String[] args) {
        String src = "43473442003802002900110500322e302e36220400303030302c01000067040004000000340500312e392e35350500322e302e361d010001140c00a82b0f6707332e00003ae6006109";
        byte[] bs = hexStringToByteArray(src);
        tlvDecode(bs);
    }
}

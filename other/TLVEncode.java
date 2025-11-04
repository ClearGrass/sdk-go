package other; 

// FIXME rename package

import java.util.*;
import java.io.ByteArrayOutputStream;

public class TLVEncoder {
    public static class Command {
        public int cmd;
        public int reportIntervl;
        public int collectInterval;
        public int valveOpen;
        public int valveSelfCheck;
        public int endFlag;
    }

    // 生成crc
    public static int byteSumU16(byte[] in) {
        int sum = 0;
        for (byte v : in) {
            sum += (v & 0xFF); // 转成无符号
        }
        return sum & 0xFFFF; // 返回 uint16 效果
    }

    public static byte[] intToBytesLittleEndian(int val, int length) {
        byte[] bytes = new byte[length];
        for (int i = 0; i < length; i++) {
            bytes[i] = (byte) (val & 0xFF);
            val >>= 8;
        }
        return bytes;
    }

    public static byte[] tlvEncode(Command cmd) {
    try {
        int cmdType = 0x32;

        ByteArrayOutputStream reportIntervlPart = new ByteArrayOutputStream();
        ByteArrayOutputStream collectIntervalPart = new ByteArrayOutputStream();
        ByteArrayOutputStream valveOpenPart = new ByteArrayOutputStream();
        ByteArrayOutputStream valveSelfCheckPart = new ByteArrayOutputStream();
        ByteArrayOutputStream endFlagPart = new ByteArrayOutputStream();

        int size = 0;
        if (cmd.collectInterval > 0) {
            int partLen = 2;
            byte[] collectBytes = intToBytesLittleEndian(cmd.collectInterval, partLen);
            byte[] partLenBytes = intToBytesLittleEndian(partLen, 2);
            collectIntervalPart.write(0x05);
            collectIntervalPart.write(partLenBytes);
            collectIntervalPart.write(collectBytes);
            size += 5;
        }
        if (cmd.reportIntervl > 0) {
            int partLen = 2;
            byte[] reportIntervlBytes = intToBytesLittleEndian(cmd.reportIntervl / 60, partLen);
            byte[] partLenBytes = intToBytesLittleEndian(partLen, 2);
            reportIntervlPart.write(0x04);
            reportIntervlPart.write(partLenBytes);
            reportIntervlPart.write(reportIntervlBytes);
            size += 5;
        }
        
        if (cmd.valveOpen > 0) {
            int partLen = 2;
            byte[] valveBytes = intToBytesLittleEndian(cmd.valveOpen*10, partLen);
            byte[] partLenBytes = intToBytesLittleEndian(partLen, 2);
            valveOpenPart.write(0x72);
            valveOpenPart.write(partLenBytes);
            valveOpenPart.write(valveBytes);
            size += 5;
            cmdType = 0x3D;
        }

        if (cmd.valveSelfCheck > 0) {
            int partLen = 1;
            byte[] partLenBytes = intToBytesLittleEndian(partLen, 2);
            valveOpenPart.write(0x73);
            valveOpenPart.write(partLenBytes);
            valveOpenPart.write(0);
            size += 4;
            cmdType = 0x3D;
        }

        if (cmd.endFlag > 0) {
            int partLen = 1;
            byte[] partLenBytes = intToBytesLittleEndian(partLen, 2);
            endFlagPart.write(0x1D);
            endFlagPart.write(partLenBytes);
            endFlagPart.write(cmd.endFlag);
            size += 4;
        }

        byte[] sizeByte = intToBytesLittleEndian(size, 2);

        ByteArrayOutputStream tlvEncode = new ByteArrayOutputStream();
        tlvEncode.write(0x43);
        tlvEncode.write(0x47);
        tlvEncode.write(cmdType);
        tlvEncode.write(sizeByte);

       
        if (collectIntervalPart.size() > 0) {
            tlvEncode.write(collectIntervalPart.toByteArray());
        }
        if (reportIntervlPart.size() > 0) {
            tlvEncode.write(reportIntervlPart.toByteArray());
        }
        if (valveOpenPart.size() > 0) {
            tlvEncode.write(valveOpenPart.toByteArray());
        }
        if (endFlagPart.size() > 0) {
            tlvEncode.write(endFlagPart.toByteArray());
        }

        int crc = byteSumU16(tlvEncode.toByteArray());
        tlvEncode.write(intToBytesLittleEndian(crc, 2));

        return tlvEncode.toByteArray();  // ✅ 成功返回
    } catch (Exception e) {
        e.printStackTrace();
        return new byte[0]; // ✅ 捕获异常时返回空数组
    }
}

    public static String bytesToHex(byte[] bytes) {
        StringBuilder sb = new StringBuilder();
        for (byte b : bytes) {
            sb.append(String.format("%02X", b)); // 两位大写十六进制
        }
        return sb.toString();
    }

    // 主方法
    public static void main(String[] args) {
        TLVEncoder encode = new TLVEncoder();
        Command cmd = new Command();
        cmd.reportIntervl = 3600;
        cmd.collectInterval = 600;
        cmd.valveOpen = 50;
        cmd.valveSelfCheck = 1;
        cmd.endFlag = 1;
        byte[] bs = encode.tlvEncode(cmd);

        // 下发命令 直接下发数组，hex编码是为方便debug
        System.out.println(bytesToHex(bs));
    }
}

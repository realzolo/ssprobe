#!/usr/bin/env python3   

PORT = 3384 
INTERVAL = 1 
PROBE_PROTOCOL_PREFER = "ipv4"  # ipv4, ipv6
# The packet loss rate monitoring direction can be customized, for example, CU = "www.baidu.com".  
CU = "www.baidu.com"
CT = "www.baidu.com"
CM = "www.baidu.com"

import socket
import time
import timeit
import platform
import os
import sys
import json
import errno
import psutil
import threading
try:
    from queue import Queue     # python3
except ImportError:
    from Queue import Queue     # python2

# Get the operating system type.
def get_platform():
    return platform.system()

# Get the operating system run time. The unit is second.
def get_uptime():
    return int(time.time() - psutil.boot_time())

# Get the usage of the memory, return total memory and used memory in MB.
def get_memory():
    Mem = psutil.virtual_memory()
    return int(Mem.total / 1024.0), int(Mem.used / 1024.0)

# Get the usage of the system swap memory, return total swap memory and used swap memory in MB.
def get_swap():
    Mem = psutil.swap_memory()
    return int(Mem.total/1024.0), int(Mem.used/1024.0)

# Get the usage of disks.
def get_hdd():
    valid_fs = [ "ext4", "ext3", "ext2", "reiserfs", "jfs", "btrfs", "fuseblk", "zfs", "simfs", "ntfs", "fat32", "exfat", "xfs" ]
    disks = dict()
    size = 0
    used = 0
    for disk in psutil.disk_partitions():
        if not disk.device in disks and disk.fstype.lower() in valid_fs:
            disks[disk.device] = disk.mountpoint
    for disk in disks.values():
        usage = psutil.disk_usage(disk)
        size += usage.total
        used += usage.used
    return int(size/1024.0/1024.0), int(used/1024.0/1024.0)

# Get the usage of CPU, return the percentage of CPU utilization.
def get_cpu():
    return psutil.cpu_percent(interval=INTERVAL)

# Get traffic information.(Unit: bytes) 
def get_traffic():
    traffic_in = 0
    traffic_out = 0
    net = psutil.net_io_counters(pernic=True)
    for k, v in net.items():
        if 'lo' in k or 'tun' in k \
                or 'docker' in k or 'veth' in k \
                or 'br-' in k or 'vmbr' in k \
                or 'vnet' in k or 'kube' in k:
            continue
        else:
            traffic_in += v[1]
            traffic_out += v[0]
    return traffic_in, traffic_out

def tupd():
    '''
    tcp, udp, process, thread count: for view ddcc attack , then send warning
    :return:
    '''
    try:
        if sys.platform.startswith("linux") is True:
            t = int(os.popen('ss -t|wc -l').read()[:-1])-1
            u = int(os.popen('ss -u|wc -l').read()[:-1])-1
            p = int(os.popen('ps -ef|wc -l').read()[:-1])-2
            d = int(os.popen('ps -eLf|wc -l').read()[:-1])-2
        elif sys.platform.startswith("win") is True:
            t = int(os.popen('netstat -an|find "TCP" /c').read()[:-1])-1
            u = int(os.popen('netstat -an|find "UDP" /c').read()[:-1])-1
            p = len(psutil.pids())
            d = 0
            # cpu is high, default: 0
            # d = sum([psutil.Process(k).num_threads() for k in [x for x in psutil.pids()]])
        else:
            t,u,p,d = 0,0,0,0
        return t,u,p,d
    except:
        return 0,0,0,0

def get_network(ip_version):
    if(ip_version == 4):
        HOST = "ipv4.google.com"
    elif(ip_version == 6):
        HOST = "ipv6.google.com"
    try:
        socket.create_connection((HOST, 80), 2).close()
        return True
    except:
        return False

lostRate = {
    '10010': 0.0,
    '189': 0.0,
    '10086': 0.0
}
pingTime = {
    '10010': 0,
    '189': 0,
    '10086': 0
}
netSpeed = {
    'netrx': 0.0,
    'nettx': 0.0,
    'clock': 0.0,
    'diff': 0.0,
    'avgrx': 0,
    'avgtx': 0
}

# Start a ping thread and continuously getting ping data.
def _ping_thread(host, mark, port):
    lostPacket = 0
    packet_queue = Queue(maxsize=100)

    IP = host
    if host.count(':') < 1:     # if not plain ipv6 address, means ipv4 address or hostname
        try:
            if PROBE_PROTOCOL_PREFER == 'ipv4':
                IP = socket.getaddrinfo(host, None, socket.AF_INET)[0][4][0]
            else:
                IP = socket.getaddrinfo(host, None, socket.AF_INET6)[0][4][0]
        except Exception:
                pass

    while True:
        if packet_queue.full():
            if packet_queue.get() == 0:
                lostPacket -= 1
        try:
            b = timeit.default_timer()
            socket.create_connection((IP, port), timeout=1).close()
            pingTime[mark] = int((timeit.default_timer() - b) * 1000)
            packet_queue.put(1)
        except socket.error as error:
            if error.errno == errno.ECONNREFUSED:
                pingTime[mark] = int((timeit.default_timer() - b) * 1000)
                packet_queue.put(1)
            #elif error.errno == errno.ETIMEDOUT:
            else:
                lostPacket += 1
                packet_queue.put(0)

        if packet_queue.qsize() > 30:
            lostRate[mark] = float(lostPacket) / packet_queue.qsize()

        time.sleep(INTERVAL)

# Continuously get network speed.
def _net_speed():
    while True:
        avgrx = 0
        avgtx = 0
        for name, stats in psutil.net_io_counters(pernic=True).items():
            if "lo" in name or "tun" in name \
                    or "docker" in name or "veth" in name \
                    or "br-" in name or "vmbr" in name \
                    or "vnet" in name or "kube" in name:
                continue
            avgrx += stats.bytes_recv
            avgtx += stats.bytes_sent
        now_clock = time.time()
        netSpeed["diff"] = now_clock - netSpeed["clock"]
        netSpeed["clock"] = now_clock
        netSpeed["netrx"] = int((avgrx - netSpeed["avgrx"]) / netSpeed["diff"])
        netSpeed["nettx"] = int((avgtx - netSpeed["avgtx"]) / netSpeed["diff"])
        netSpeed["avgrx"] = avgrx
        netSpeed["avgtx"] = avgtx
        time.sleep(INTERVAL)

# Get real-time data(ping data and net speed).
def get_realtime_data():
    t1 = threading.Thread(
        target=_ping_thread,
        kwargs={
            'host': CU,
            'mark': '10010',
            'port': 80
        }
    )
    t2 = threading.Thread(
        target=_ping_thread,
        kwargs={
            'host': CT,
            'mark': '189',
            'port': 80
        }
    )
    t3 = threading.Thread(
        target=_ping_thread,
        kwargs={
            'host': CM,
            'mark': '10086',
            'port': 80
        }
    )
    t4 = threading.Thread(
        target=_net_speed,
    )
    t1.deamon = True
    t2.deamon = True
    t3.deamon = True
    t4.deamon = True
    t1.start()
    t2.start()
    t3.start()
    t4.start()

# String and byte conversion.
def byte_str(object):
    if isinstance(object, str):
        return object.encode(encoding="utf-8")
    elif isinstance(object, bytes):
        return bytes.decode(object)
    else:
        print(type(object))

# User host authentication.
def auth(skt,token):
    skt.send(byte_str(token))
    authResObj = json.loads(byte_str(skt.recv(1024)))
    if authResObj['code'] == -1: # auth_res: -1(认证失败) 0(认证成功)
        print("Authentication failure!")
        return False,None
    else:
        print("Authentication succeeded, you have successfully connected to the server!")
        return True,authResObj

# Get the initialization data entered by the user.
def getInitData():
    server = ""
    name = ""
    token = ""
    while 1:
        server = input("Server address: ")
        name = input("Client name: ")
        token = input("Token: ")
        if(len(server)==0 or len(name) == 0 or len(token) ==0):
            print("The content cannot be empty!")
        else:
            break
    return server,name,token
    

if __name__ == '__main__':
    # 获取用户初始化信息
    server,name,token = getInitData()
    socket.setdefaulttimeout(30)
    get_realtime_data()
    while 1:
        try:
            print("Connecting...")
            skt = socket.create_connection((server, PORT))

            # 用户认证
            ok,authResObj = auth(skt,token)
            if not ok:
                os._exit(0)
            
            timer = 0
            check_ip = 0
            if authResObj['ip_version'] == 'IPv4':
                check_ip = 4
            elif authResObj['ip_version'] == 'IPv6':
                check_ip = 6

            while 1:
                CPU = get_cpu()
                TrafficIn, TrafficOut = get_traffic()
                Platform = get_platform()
                Uptime = get_uptime()
                Load_1, Load_5, Load_15 = os.getloadavg() if 'linux' in sys.platform else (0.0, 0.0, 0.0)
                MemoryTotal, MemoryUsed = get_memory()
                SwapTotal, SwapUsed = get_swap()
                HDDTotal, HDDUsed = get_hdd()

                array = {}
                if not timer:
                    array['online' + str(check_ip)] = get_network(check_ip)
                    timer = 10
                else:
                    timer -= 1*INTERVAL

                array['name'] = name
                array['host'] = authResObj['host']
                array['state'] = True
                array['platform'] = Platform
                array['uptime'] = Uptime
                array['load_1'] = Load_1
                array['load_5'] = Load_5
                array['load_15'] = Load_15
                array['memory_total'] = MemoryTotal
                array['memory_used'] = MemoryUsed
                array['swap_total'] = SwapTotal
                array['swap_used'] = SwapUsed
                array['hdd_total'] = HDDTotal
                array['hdd_used'] = HDDUsed
                array['cpu'] = CPU
                array['network_rx'] = netSpeed.get("netrx")
                array['network_tx'] = netSpeed.get("nettx")
                array['network_in'] = 1 << 30
                array['network_out'] = TrafficOut
                array['ping_10010'] = lostRate.get('10010') * 100
                array['ping_189'] = lostRate.get('189') * 100
                array['ping_10086'] = lostRate.get('10086') * 100
                array['time_10010'] = pingTime.get('10010')
                array['time_189'] = pingTime.get('189')
                array['time_10086'] = pingTime.get('10086')
                array['tcp'], array['udp'], array['process'], array['thread'] = tupd()

                skt.send(byte_str(json.dumps(array)))
        except KeyboardInterrupt:
            print("Disconnected!")
            if 'skt' in locals().keys():
                del skt
            os._exit(0)        
        except socket.error:
            print("Disconnected!")
            if 'skt' in locals().keys():
                del skt
            time.sleep(3)
        except Exception as e:
            print("Caught Exception:", e)
            if 'skt' in locals().keys():
                del skt
            time.sleep(3)

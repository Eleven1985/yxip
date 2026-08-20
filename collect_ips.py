import requests
from bs4 import BeautifulSoup
import re
import os
import ssl
from requests.adapters import HTTPAdapter
from urllib3.poolmanager import PoolManager

# 自定义 TLS 适配器，解决 SSL 握手失败问题
class TLSAdapter(HTTPAdapter):
    def init_poolmanager(self, *args, **kwargs):
        context = ssl.create_default_context()  # 默认启用 TLS1.2+
        kwargs['ssl_context'] = context
        return super().init_poolmanager(*args, **kwargs)

# 设置目标 URL 列表
urls = [
    'https://monitor.gacjie.cn/page/cloudflare/ipv4.html',
    'https://ip.164746.xyz'
]

# IP 地址正则表达式
ip_pattern = r'\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}'

# 如果 ip.txt 存在则删除
if os.path.exists('ip.txt'):
    os.remove('ip.txt')

# 创建 requests session 并挂载自定义 TLSAdapter
session = requests.Session()
session.mount('https://', TLSAdapter())

# 写入文件
with open('ip.txt', 'w') as file:
    for url in urls:
        try:
            response = session.get(url, timeout=10)
            response.raise_for_status()
        except requests.exceptions.RequestException as e:
            print(f"[错误] 无法请求 {url}：{e}")
            continue

        soup = BeautifulSoup(response.text, 'html.parser')

        # 根据网址结构获取 IP 所在元素
        if url in ['https://monitor.gacjie.cn/page/cloudflare/ipv4.html', 'https://ip.164746.xyz']:
            elements = soup.find_all('tr')
        else:
            elements = soup.find_all('li')

        # 提取 IP 地址
        for element in elements:
            element_text = element.get_text()
            ip_matches = re.findall(ip_pattern, element_text)
            for ip in ip_matches:
                file.write(ip + '\n')

print("✅ 所有 IP 地址已保存到 ip.txt 文件中。")

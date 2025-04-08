# nmap_tool.yaml
# ======================================
# 工具元数据
# ======================================
name: nmap  # 工具唯一标识
label: Nmap 端口扫描器
description: 开源网络探测和安全审计工具，用于扫描开放端口和服务识别
version: 7.94

# ======================================
# 输入输出定义
# ======================================
input_types:  # 支持的输入类型（可多选）
  - ipv4
  - ipv6
  - domain
output_type: port_scan_result  # 全局唯一输出类型标识

# ======================================
# 参数定义（支持动态生成命令行）
# ======================================
parameters:
  # 扫描类型选择（枚举值）
  - name: scan_type
    type: enum
    values: 
      - tcp_syn    # SYN半开扫描
      - tcp_connect  # 全连接扫描
      - udp
      - quick       # 快速扫描
    default: tcp_syn
    cli_arg: "-s%s"  # 参数映射规则（%s自动替换值）
    description: 选择扫描技术类型

  # 操作系统检测
  - name: os_detection
    type: boolean
    default: false
    cli_arg: "-O"
    description: 启用操作系统指纹识别

  # 端口范围配置
  - name: ports
    type: string
    default: "1-1024"
    cli_arg: "-p %s"
    validation: "^\\d+(-\\d+)?(,\\d+(-\\d+)?)*$"  # 正则校验格式
    description: 指定扫描端口范围（示例：80,443,1000-2000）

  # 输出格式控制（影响后续解析逻辑）
  - name: output_format
    type: enum
    values: [xml, json]
    default: xml
    cli_arg: "-o%s"
    description: 指定输出格式（强烈建议使用XML/JSON以便解析）

# ======================================
# 执行配置（Docker模式）
# ======================================
execution:
  type: docker
  image: instrumentisto/nmap:7.94  # 官方推荐镜像
  command_template: |  # 动态生成命令
    nmap ${scan_type} ${os_detection} 
    -p ${ports} ${output_format} 
    -n -T4 ${input_target}

  resource_limits:  # 资源限制
    memory: 512M
    cpu: 1

# ======================================
# 输出解析规则（关键！）
# ======================================
output_parser:
  xml:  # 当output_format=xml时使用的解析器
    type: xpath
    schema:
      ip: /nmaprun/host/address/@addr
      ports: 
        path: /nmaprun/host/ports/port
        attributes:
          port_id: @portid
          protocol: @protocol
          service: service/@name
          state: state/@state
      os: /nmaprun/host/os/osmatch/@name
  
  json:  # 当output_format=json时使用的解析器
    type: jsonpath
    schema:
      ip: $.host.$.address.$.addr
      ports: 
        path: $.host.$.ports.$.port
        fields:
          port: $.portid
          protocol: $.protocol
          service: $.service.name
          state: $.state.state

# ======================================
# 测试用例（用于AI学习和验证）
# ======================================
test_cases:
  - name: 基础TCP扫描
    input: 
      target: "192.168.1.1"
      params:
        scan_type: tcp_syn
        ports: "80,443"
    expected_output:
      ip: "192.168.1.1"
      ports:
        - port: 80
          protocol: tcp
          state: open
          service: http
        - port: 443
          protocol: tcp  
          state: open
          service: https

  - name: 完整扫描示例
    input: 
      target: "example.com"
      params:
        scan_type: tcp_connect
        ports: "1-10000"
        os_detection: true
        output_format: xml
    expected_output:
      ip: "93.184.216.34"  # example.com的IP
      os: "Linux 3.10+"
      ports: 
        - port: 80
          protocol: tcp
          state: open
          service: http
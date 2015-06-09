// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// 基于xml配置的日志系统。
//
// logs 定义了6个级别的日志：ERROR,INFO,TRACE,DEBUG,CRITICAL,WARN。
// 用户可以根据自己的需求，通过xml配置文件自定义每个日志输出行为。
// 以下为一个简短的xml配置示例，具体的可参考目录下的config.xml。
//  <?xml version="1.0" encoding="utf-8" ?>
//  <logs>
//      <debug>
//          <buffer size="10">
//              <rotate dir="/var/log/" size="5M" />
//              <stmp username=".." password=".." />
//          </buffer>
//      </debug>
//      <info>
//          <console output="stderr" color="yellow" />
//      </info>
//      <!-- 除了debug和info，其它4个依然输出到ioutil.Discard -->
//  </logs>
//
// 然后就可以调用go代码输出日志内容:
//  logs.Debug(...)
//  logs.Debugf("format", v...)
//  logs.DEBUG.Println(...)
//
//  // error并未在配置文件中出现，所有内容将不作实际输出()
//  logs.ERROR().Print(...)
//
//  // 向所有级别的日志输出内容。
//  logs.All(...)
//
// 上面的配置文件表示DEBUG级别的内容输出前都进被缓存，当量达到10条时，
// 一次性向rotate和stmp输出。
// 其中buffer、rotate、stmp甚至是debug和info都是一个个实现io.Writer
// 接口的结构。通过Register()注册成功之后，即可以使用。
//
//
//
// 配置文件：
//
// - 只支持utf-8编码的xml文件。
//
// - 节点名称和节点属性区分大小写，但是属性值不区分大小写。
//
// - 顶级元素必须为logs，且不需要带任何属性;
//
// - 二级元素只能为info,deubg,trace,warn,error,critical。
// 分别对应INFO,DEBUG,TRACE,WARN,ERROR,CRITICAL等日志实例。
// 可以带上prefix和flag属性，分别对应log.New()中的相应参数。
//
// - 三级及以下元素可以自己根据需求组合，logs自带以下writer，
// 用户也可以自己向logs注册自己的实现。
//
// 1. buffer:
//
// 缓存工具，当数量达到指定值时，一起向所有的子元素输出。
// 比如上面的示例中，所有向debug输出的内容，都会被buffer缓存，
// 直到数量达到10条，才会一起向rotate和stmp输出内容。
// 仅有size一个参数：
//  size: 用于指定缓存的数量，必填参数。
//
// 2. rotate:
//
// 这是一个按文件大写自动分割日志的实例，以第一条记录的产生时间作为文件名。
// 拥有以下三个参数：
//  prefix：表示日志文件的前缀，留空表示没有前缀；
//  dir：	表示的是日志存放的目录；
//  size：	表示的是每个日志的大概大小，默认单位为byte，可以带字符单位，
//          如5M,10G等(支持k,m,g三个后缀，不区分大小写)。
//
// 3. stmp:
//
// 将日志内容发送给指定邮件，可定义的属性为：
//  username: 发送邮件的账号；
//  password: 账号对应的密码；
//  host:	  stmp的主机；
//  subject:  邮件的主题；
//  sendTo:   接收人地址，多个收件地址使用分号分隔。
//
// 4. console:
//
// 向控制台输出内容。可定义的属性为：
//  output：    只能为"stderr", "stdout"两个值，表示输出的具体方向，默认值为"stderr"；
//  foreground: 表示输出时的前景色，其值在github.com/issue9/term/colors中定义。
//  background: 表示输出时的背景色，其值在github.com/issue9/term/colors中定义。
//
//
// 自定义
//
// 除了以上定义的元素，用户也可以自行实现io.Writer接口，以实现自定义的输出方向。
// 要添加新的元素，只需要向Register()函数注册一个元素的初始化函数即可，
// 其中注册的名称将作为配置节点的元素名称，需要唯一，不能与现有的名称相同；
// 函数则作为节点转换成实例时的转换功能，需要负责解析节点传递过来的属性列表(args参数)，
// 若是一个容器节点（如buffer，可以包含子节点）则返回的实例必须要实现WriteFlushAdder接口，
// 该函数的原型为：
//  WriterInitializer
package logs

const Version = "0.8.15.150609"

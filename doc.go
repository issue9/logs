// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// logs 是对标准库的log的一个简单扩展，定义了6个级别的日志：
// ERROR,INFO,TRACE,DEBUG,CRITICAL,WARN，用户可以根据自己的需求，
// 向指定级别的日志输出特定内容，也可以通过xml文件定义每个级别的日志行为。
//
// 默认情况下，所有日志都输出到ioutil.Discard。
// 用户需要显示地调用相应的初始化函数对日志进行初始化。
//
// 以下是一个简短的xml配置文件范本，具体的可参考目录下的config.xml。
//  xml:
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
//  go:
//  logs.Debug(...)
//  logs.Debugf("format", v...)
//  logs.DEBUG.Println(...)
//
//  // error并未在配置文件中出现，所有内容将输出到ioutil.Discard
//  logs.ERROR.Print(...)
//
//  // 向所有级别的日志输出内容。
//  logs.All(...)
//
// 上面的配置文件表示DEBUG级别的内容输出前都进被buffer实例进行缓存，
// 当量达到10条时，一次性向rotate和stmp输出。
// 其中buffer、rotate、stmp甚至是debug和info都是一个个实现io.Writer
// 接口的结构。通过Register()注册成功之后，即可以使用。
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
// - 三级元素可以自己根据需求组合，logs自带以下writer，
// 用户也可以自己向logs注册自己的实现。
//
// 1. buffer:
// 缓存工具，比如上面的示例中，所有向debug输出的内容，都会被buffer缓存，
// 直到数量达到10条，buffer才会一起向rotate和stmp输出内容。
// 只有一个size属性，用于指定缓存的数量。
//
// 2. rotate:
// 这是一个按文件大写自动分割日志的实例，允许dir和size两个属性。
// 其中dir表示的是日志存放的目录；size表示的是每个日志的大概大小，
// 可以是数值(单位为Byte)或是5M,5G等这类字符串(支持k,m,g三个后缀，不区分大小写)。
//
// 3. stmp:
// 发送邮件的实例，可定义的属性为：username 发送邮件的账号；
// password账号对应的密码；host stmp的主机；subject 邮件的主题；sendTo，
// 接收人地址，多个收件地址使用分号分隔。
//
// 4. console:
// 向控制台输出内容。可定义的属性为：foreground，background和output，
// 其中output只能为"stderr", "stdout"两个值，表示输出的具体方向，
// 默认值为"stderr"；而foreground和background表示输出时的前景和背景色，
// 其值在github.com/issue9/term/colors中定义。
//
// 自定义
//
// 除了以上定义的元素，用户也可以自定义一些新io.Writer以实现自定义的输出功能。
// 要添加新的元素，只需要向Register()函数注册一个元素的初始化函数即可，
// 其中注册的名称将作为配置节点的元素名称，需要唯一，不能与现有的名称相同；
// 函数则作为节点转换成实例时的转换功能，需要负责解析节点传递过来的属性列表(args参数)，
// 若是一个容器节点（如buffer，可以包含子节点）则返回的实例必须要实现WriteFlushAdder接口，
// 该函数的原型为：
//  func(args map[string]string) (io.Writer, error)
// 具体可参考WriterInitializer。
package logs

const Version = "0.2.7.150121"

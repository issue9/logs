// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// logs为日志处理包，相对于标准库的log，添加了对分级日志的支持，
// 以及在不需要重新编译源代码的情况下，修改日志输出行为的功能。
//
// 以下是一个简短的xml配置文件范本，具体的可参考目录下的config.xml
// 文件。
//  <?xml version="1.0" encoding="utf-8" ?>
//  <logs>
//      <debug>
//          <buffer size="10">
//              <rotate dir="/var/log/" size="5M" />
//              <stmp username=".." password=".." />
//          </buffer>
//      </debug>
//      <info>
//          ....
//      </info>
//  </logs>
// 上面的配置文件表示LevelDebug级别的内容输出前都进被buffer实例进行
// 缓存，当量达到10条时，一次性向rotate和stmp输出。
// 其中buffer、rotate、stmp甚至是debug和logs都是一个个实现io.Writer
// 接口的结构。通过Register()向logs注册成功之后，即可以使用。
//
//
//
// 以下是包自带io.Writer实例，可以直接使用：
//
// 1. logs:
// 顶级元素，包会自动替换成LevelLogger。
//
// 2. info,deubg,trace,warn,error,critical:
// 配置文件中的二级元素，不支持自定义，填写其它名称会返回error。可以
// 指定prfix和flag属性，分别对应log.New()的prefix和flag参数，事实上，
// 最终也是调用log.Logger结构来输出日志的。
//
// 3. buffer:
// 一个简单的缓存工具，比如上面的示例中，所有向debug输出的内容，都会
// 被buffer缓存，直到数量达到10条，buffer才会一起向rotate和stmp输出内
// 容。只有一个size属性，用于指定缓存的数量。
//
// 4. rotate:
// 这是一个按文件大写自动分割日志的实例，允许dir和size两个属性。其中
// dir表示的是日志存放的目录；size表示的是每个日志的大概大小，可以是
// 数值(单位为Byte)或是5M,5G等这类字符串(具体哪类字符串会被正确解析，
// 可参考conv包的ToByte()函数)。
//
// 5. stmp:
// 发送邮件的实例，可定义的属性为：username 发送邮件的账号；password
// 账号对应的密码；host stmp的主机；subject 邮件的主题；sendTo，接收
// 人地址，多个收件地址使用分号分隔。
//
// 6. console:
// 向控制台输出内容。可定义的属性为：color和output，其中output只能为
// "os.Stderr","os.stdin", "os.stdout"三个值，表示输出的具体方向，默
// 认值为"os.stderr"；而color表示输出时的颜色，为一个ansi颜色控制码
// (在windows将原样输出),不指定时，默认为红包，term包中定义了部分颜色
// 值，可以直接拿来用。
//
//
// 其它一些注意事项：
//
// 1. 配置文件只支持utf-8编码；
// 2. 元素和属性区分大小写，但属性值不区分大小写；
// 3. 根元素必须为logs；
// 4. 二级元素只能是info、debug、trace、warn、error和critical分别对应
// 各个level元素；
//
//
//
// 自定义
//
// 除了以上定义的元素，用户也可以自定义一些新io.Writer以实现自定义的输
// 出功能。要添加新的元素，只需要向Register()函数注册一个元素的初始化函
// 数即可，其中注册的名称将作为配置节点的元素名称，需要唯一，不能与现有
// 的名称相同；函数则作为节点转换成实例时的转换功能，需要负责解析节点传
// 递过来的属性列表(args参数)，若是一个容器节点（如buffer，可以包含子节
// 点）则返回的实例必须要实现WriteFlushAdder接口，该函数的原型为：
//  func(args map[string]string) (io.Writer, error)
// 具体可参考WriterInitializer。
//
// 注意事项：
//
// 根元素和二级元素不能自定义，也不能将根元素或是二级元素用于其它地方。
package logs

const Version = "0.1.7.141023"

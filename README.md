基于epoll的事件中心，支持建立事件组：
1.单个上下文中监听多个不同的事件通知，进行事件处理；
2.支持动态的增加事件或删除事件；
3.支持暂停某个事件和恢复某个事件的通知。
// pages/detail/detail.js

var WxParse = require('../../wxParse/wxParse.js');

Page({
  /**
   * 页面的初始数据
   */
  data: {
    contentTip: '由于后台接口原因，新闻具体内容无法编辑，只返回了一个新闻链接...'
  },


  /**
   * 生命周期函数--监听页面加载
   */
  onLoad: function (options) {
    var that = this;
    wx.request({
      url: `http://152.136.128.83:8080/web_news/news/detail/` + options.newsId,
      
      success(res) {
        
        var article = res.data.text;
        WxParse.wxParse('article', 'html', article, that, 5);

        
      }
    })
  },
})
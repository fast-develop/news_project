//
//  DetailViewController.h
//  ZNews
//
//  Created by Frank Zheng on 10/20/14.
//  Copyright (c) 2014 xzheng. All rights reserved.
//

#import <UIKit/UIKit.h>
#import "MOArticle.h"

@interface NewsViewController : UIViewController
/*
 {
 
 UIWebView *webView;
 }
 */
@property(nonatomic, strong) UIWebView *webView;
@property (strong, nonatomic) MOArticle *detailItem;

@end


//
//  DetailViewController.m
//  ZNews
//
//  Created by Frank Zheng on 10/20/14.
//  Copyright (c) 2014 xzheng. All rights reserved.
//

#import "NewsViewController.h"
#import "ContentService.h"
#import "MOArticleDetail+Dao.h"
#import "ModelUtil.h"

#define SCREEN_WIDTH [UIScreen mainScreen].bounds.size.width //屏幕宽度

#define SCREEN_HEIGHT [UIScreen mainScreen].bounds.size.height //屏幕高度

@interface NewsViewController ()
@end

@implementation NewsViewController

- (void)setDetailItem:(MOArticle *)newDetailItem {
    
    if (_detailItem != newDetailItem) {
        _detailItem = newDetailItem;
    }
    
    //[self configView];
}


- (void)viewDidLoad
{
    /*
     [super viewDidLoad];
     webView = [[UIWebView alloc] initWithFrame:CGRectMake(0, 0, 320, 480)];
     NSURLRequest *request =[NSURLRequest requestWithURL:[NSURL URLWithString:@"http://www.baidu.com"]];
     [self.view addSubview: webView];
     [webView loadRequest:request];
     */
    
    [super viewDidLoad];
    self.webView = [[UIWebView alloc] initWithFrame:CGRectMake(0, 0, SCREEN_WIDTH, SCREEN_HEIGHT)];
    //NSURLRequest *request =[NSURLRequest requestWithURL:[NSURL URLWithString:@"http://www.baidu.com"]];
    [self.view addSubview: self.webView];
    
    NSString *CSS= @"<style type=\"text/css\">img{ width:100%;}</style>";
    
    static NSString *body_new =@"";
    if(self.detailItem.detail == nil)
    {
        body_new = @"Loading...";
        
        //load article detail from backend
        [[ContentService instance] getArticleDetail:self.detailItem sucess:^(NSDictionary *data) {
            //insert article detail to db
            MOArticleDetail *detail = [MOArticleDetail insertArticleDetailWithDictionary:data
                inManagedObjectContext:defaultManagedObjectContext()
                relatedToArticle:self.detailItem];
            
            //save the changes
            commitDefaultMOC();
            //NSLog(detail.text);
            
            body_new = [detail.text mutableCopy];
            
            //NSLog(body_new);
            
            NSString * htmlString = [NSString stringWithFormat:@"<html><meta charset=\"UTF-8\"><header>%@</header><body>%@</body></html>",CSS,body_new];
            [self.webView loadHTMLString:htmlString baseURL:nil];
            
        } failure:^{
            body_new = @"Unable to load the article text.";
            NSString * htmlString = [NSString stringWithFormat:@"<html><meta charset=\"UTF-8\"><header>%@</header><body>%@</body></html>",CSS,body_new];
            [self.webView loadHTMLString:htmlString baseURL:nil];
        }];
    }
    else
    {
        body_new = [self.detailItem.detail.text mutableCopy];
        NSString * htmlString = [NSString stringWithFormat:@"<html><meta charset=\"UTF-8\"><header>%@</header><body>%@</body></html>",CSS,body_new];
        [self.webView loadHTMLString:htmlString baseURL:nil];
    }
    
    
    //NSLog(body_new);
    NSString * htmlString = [NSString stringWithFormat:@"<html><meta charset=\"UTF-8\"><header>%@</header><body>%@</body></html>",CSS,body_new];
    [self.webView loadHTMLString:htmlString baseURL:nil];
    
}

@end


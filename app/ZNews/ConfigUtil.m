//
//  ConfigUtil.m
//  ZNews
//
//  Created by xuyi on 2019/6/4.
//  Copyright © 2019 xzheng. All rights reserved.
//

#import "ConfigUtil.h"

@implementation ConfigUtil

+ (ConfigUtil*)getInstance
{
    static ConfigUtil* instance = nil;
    static dispatch_once_t once;
    
    dispatch_once(&once, ^{
        instance = [[self.class alloc] init];
    });
    return instance;
}

-(NSString*)valueOf:(NSString*)name
{
    NSString * serverUrl = nil;
    NSString * filePath= [[NSBundle mainBundle] pathForResource:@"Settings" ofType:@"plist"];
    NSDictionary * contentsOfFile =  [[NSDictionary alloc]initWithContentsOfFile:filePath];
    return [contentsOfFile objectForKey:name];
}

@end

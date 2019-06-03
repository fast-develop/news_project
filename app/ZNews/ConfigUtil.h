//
//  ConfigUtil.h
//  ZNews
//
//  Created by xuyi on 2019/6/4.
//  Copyright © 2019 xzheng. All rights reserved.
//

#import <Foundation/Foundation.h>

@interface ConfigUtil : NSObject

+ (ConfigUtil*)getInstance;
-(NSString*)valueOf:(NSString*)name;

@end

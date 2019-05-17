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
    
    /*
     NSString *body_new = @"<div><p>本文意在了解为什么当今大多数企业中使用混合,以及此计算模型与仅使用云和无云方案的比较。</p><p>混合IT是一种企业计算模型，其中企业或组织通过传统的内部IT系统提供一些资源，同时还将云计算服务的某种组合用于其他资源。</p><p>正如其普遍存在的现象那样，对于大多数组织而言，混合IT仍然是一种比完全无云或完全依赖公有云和/或私有云的更可行的方法。</p><p>混合IT云架构 - 无论是公有还是私有 – 其特点是基于需求的可扩展性，用户可以轻松配置以及计量使用。</p><p>许多提供商已经在SaaS（软件即服务），IaaS（基础架构即服务）和PaaS（平台即服务）等领域采用公有云服务，从而减轻了客户在公司内部安装和管理类似资源的需求负担。混合IT也可以专注于私有云。与多租户公有云不同，私有云使用单租户架构。私有云通常在内部数据中心运行，尽管私有云托管服务肯定是可用的。</p><p>即使公有云服务激增，大多数公司都继续将至少某些资源保留在内部，不受任何云环境的影响。这是由于安全和隐私问题，或是由于将复杂遗留系统迁移到云架构的技术挑战。</p><div class=\"pgc-img\"><img src=\"http://p3.pstatp.com/large/pgc-image/8e1e047348474499beef3a4aa51e2c04\" img_width=\"730\" img_height=\"473\" alt=\"什么是混合IT云架构？\" inline=\"0\"><p class=\"pgc-img-caption\"></p></div><p><strong>一、混合IT解决方案</strong></p><p>在当今最常见的混合IT解决方案类型中，企业继续在传统IT环境中运行和维护关键IT资源，或是在本地又或是托管，同时又使用来自多个公有云提供商的其他资源。</p><p>然而，一些客户将混合IT仅仅作为临时解决方案，着眼于最终将所有IT资源迁移到云。</p><p>有三种类型的公有云解决方案可以与混合IT模型集成：</p><p><strong>SaaS</strong></p><p>• 当今最大的公有云市场，SaaS使用web交付由第三方供应商运营的应用程序。企业的最终用户可以直接从Web浏览器访问大多数SaaS应用程序。通常以每月订阅或按使用付费的方式提供，SaaS消除了在企业服务器和PC上安装和管理应用程序的需要，以及相关的CapEx硬件和软件许可证支出。</p><p>• 主流的SaaS公司提供办公套件，电子邮件和协作，CRM，HR，会计和ERP等软件。Google Apps，Cisco WebEx，Salesforce.com，SAP One和Citrix GoToMeeting是使用最广泛的SaaS应用程序中的一小部分。</p><p><strong>IaaS</strong></p><p>• IaaS的特点是用于访问，监控和管理远程基础架构的自助服务。云提供商托管包含服务器，存储和网络硬件，以及虚拟化或hypervisor层的基础架构组件，传统上存在于本地数据中心。</p><p>• 通常，提供商还提供监控，计费，日志访问，安全性，负载平衡和群集以及数据备份，复制和恢复等服务。</p><p>• 在公有云IaaS中，提供商根据资源消耗对托管和其他服务收费。客户可以通过广域网（通常是Internet）登录服务，以便进行故障排除应用程序和管理灾难恢复等用途。此外，IaaS服务持续受策略驱动，允许用户实现更高水平的自动化和编排。</p><p>• 许多初创公司租用或租赁IaaS服务，而不是购买自己的数据中心硬件和软件</p><p>• 然而，一些历史悠久的企业在其混合IT模型中包含IaaS。例如，IaaS用于扩展内部数据中心，以便在高需求时段（如圣诞购物季）处理临时工作负载。</p><p>• 尽管IaaS通常在公有云中运行，但解决方案也可作为私有云在公司自己的数据中心内使用。</p><p>• 虽然IaaS市场正在围绕亚马逊网络服务（AWS）和Microsoft Azure快速整合，但也存在其他提供商，例如包括IBM Cloud，Google Compute Engine（GCE），DigitalOcean，CenturyLink Cloud，Joyent和Rackspace Managed Cloud。</p><p><strong>PaaS</strong></p><p>• PaaS以IaaS为基础，提供软件组件中间件框架，IT组织可以使用这些框架构建新应用程序或自定义现有应用程序。面向敏捷DevOps，PaaS旨在使开发，测试和部署客户应用程序变得更容易，更具成本效益。</p><p>• 服务通常通过混合IT模型提供，该模型同时使用公有云IaaS和内部IT资源，但其他交付方法包括组合私有云和公有云的混合云模型以及仅使用私有云的私有PaaS。在所有“即服务”产品中，PaaS最快速地超越混合IT。根据Gartner最近的一份报告，已经有近一半的PaaS服务只是云计算。</p><p>• 在高度分散的PaaS领域，提供商包括Amazon Electronic Beanstalk，Apprenda，Cloud Foundry，Google App Engine（GAE），SAP Cloud Platform和Software AG Cloud等等。</p><p><strong>二、混合IT基础设施</strong></p><p>• 在仅支持SaaS的混合IT实施中，一些应用程序被移至公有云，但基础设施和软件开发组件仍保留在传统的内部IT环境中。</p><p>• 在涉及IaaS和/或PaaS的混合IT实施中，一些基础架构和/或开发组件迁移到公有云，而其他IT资源则安装在传统的内部IT环境中。但是，可以在内部访问和管理公有云资源。</p><p><strong>三、混合IT案例</strong></p><p>• 制造商可能会在内部运行高度个性化的大型机生产业务应用程序，以避免重写应用程序的复杂性，同时从云中采用SaaS应用程序来替换ERP，HR和CRM等领域中过时的传统客户端 - 服务器应用程序。</p><p>• 银行可以通过继续在传统的内部部署数据中心保护客户信息来遵守监管要求，同时还启动混合云生态系统来托管创新的金融科技创业合作伙伴。</p><p>• 出于竞争的原因，州际卡车运输公司可能会选择在托管的私有云中运行高度专有的应用程序，同时通过整合来自公开资源的气象和地图提要来更新数据</p><p>• 大型零售商可能会在内部部署的数据中心保留有价值的客户信息，同时通过利用公有云中强大的数据分析工具来处理数据，从而节省资金。</p><p>• 在iOT应用中，电力公司可以将其计费系统保留在内部部署数据中心，同时使用公有云服务作为端点从客户的仪表收集遥测数据。</p><p><strong>四、混合IT成本</strong></p><p>按需付费定价是采用混合IT的主要推动力。通过利用公有云，企业可以根据需要租用软件和硬件，而不是长期承诺购买所有的企业商业软件许可证和底层硬件。现有的内部IT人员可以从应用程序和网络管理重新分配到其他任务中去。</p><p>尽管如此，通过在内部运行和管理其他IT资源，组织仍需要继续购买一些软件许可证和硬件升级，以及继续为内部IT管理人员和任何适用的托管费用付费。</p><p>云服务提供商的计费通常包括企业未预料到的隐形成本。当涉及多个公有云提供商时，随之而来的是云计算蔓延。云成本和使用可能会隐藏在各个提供商和多个内部帐户，团队和业务部门的月度结算明细中。</p><p>因此，组织需要在提交公有云服务之前仔细审查计费协议，并实时监控云计费以识别和管理违规行为。</p><p>虽然公有云与私有云的成本效益构成了一个持续争论的话题，但可以肯定地说私有云通常涉及大量的咨询费和管理成本。另一方面，私有云成本往往更加透明。然后使用公有云仍然更为典型，但混合IT包含私有云和公有云。</p><p><strong>五、混合IT的优势和劣势</strong></p><p><strong>优势：</strong></p><p>• 组织可以继续对敏感资源施加控制、安全、隐私和法规遵从性，方法是将这些资源保留在内部。</p><p>• 同时，可以将不太敏感的资源迁移到灵活且易于扩展的云环境中。</p><p>• 混合IT环境将内部IT系统的固定成本与云服务的可变成本结合在一起，从而可以围绕IT支出进行更灵活的规划。</p><p><strong>劣势：</strong></p><p>• 许多混合IT环境仍然“偶然”发生，没有根据工作负载要求选择资源类型的战略计划。</p><p>• 目前用于混合IT云部署的典型SaaS应用程序不能像自行开发的内部应用程序那样可定制。</p><p>• 在混合IT模型下，管理和集成可能比完全云化或完全无云的环境更具挑战性。</p><p>原文链接：</p><p>https://www.datamation.com/cloud-computing/what-is-hybrid-it.html</p><blockquote><p>译者介绍</p></blockquote><p>武楠，云技术社区翻者，就职于华讯网络顾问工程师，主攻方向云计算及网络。</p></div>";
     */
    
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


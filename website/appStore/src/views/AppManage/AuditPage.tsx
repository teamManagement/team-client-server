/* eslint-disable jsx-a11y/anchor-is-valid */
import { useState, useCallback, useRef, useEffect } from 'react';
import { Button, Image, Tooltip, Spin, Tabs, Avatar } from 'antd'
import { ExclamationCircleOutlined } from '@ant-design/icons';
import png1 from '../../imgs/1.png';
import png2 from '../../imgs/2.png';
import png3 from '../../imgs/3.png';
import png4 from '../../imgs/4.png';
import png5 from '../../imgs/5.png';
import png6 from '../../imgs/6.png';
import png7 from '../../imgs/7.png';
import png8 from '../../imgs/8.png';
import png9 from '../../imgs/9.png';
import pngLogo from '../../imgs/markdown.png';

import './AuditPage.less'
import AuditRefuseModal, { AuditRefuseModalActionType } from './components/AuditRefuseModal';

enum TabKey {
    base = 'base',
    info = 'info',
    contributors = 'contributors',
    version = 'version'
}

interface IAuditPageProps {
    appId?: any;
    onCompleted?: () => void;
}

const AuditPage: React.FC<IAuditPageProps> = (props: IAuditPageProps) => {

    const [loading, setLoading] = useState<boolean>(false);
    const [appInfo, setAppInfo] = useState<any>();
    const [activeTabKey, setActiveTabKey] = useState<TabKey>(TabKey.base);
    const [aduitLoading, setAduitLoaing] = useState<boolean>();

    const fetchData = useCallback(() => {
        setLoading(true);

        var id = props.appId;
        console.info(id);
        //调用查询详情接口todo

        setLoading(false);
        setAppInfo({});
    }, [props.appId]);

    useEffect(() => {
        fetchData();
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, []);

    const refuseRef = useRef<AuditRefuseModalActionType>();

    const getTabItem = () => {
        let rst = [
            { label: '基础信息', key: TabKey.base },
            { label: '详细介绍', key: TabKey.info },
            { label: '贡献者列表', key: TabKey.contributors },
        ];
        if (true) {
            rst.push({ label: '版本记录', key: TabKey.version });
        }
        return rst;
    };

    const testPicList = [{ src: png1 }, { src: png2 }, { src: png3 }, { src: png4 }, { src: png5 }, { src: png6 }, { src: png7 }, { src: png8 }, { src: png9 }];
    const testContributors = [
        { userNick: 'teamworder', userName: '程亮', roleName: '超级管理员' },
        { userNick: 'teamworder', userName: '程亮亮', roleName: '超级管理员' },
    ];
    const testVersions = [
        { title: '这里是标题', describe: '这里是描述，这里是描述，这里是描述，这里是描述' },
        { title: '这里是标题', describe: '这里是描述，这里是描述，这里是描述，这里是描述' },
    ]
    const getUserNameLastDigit = (userName: string) => {
        if (userName.length <= 2) {
            return userName;
        }
        return `${userName[userName.length - 2]}${userName[userName.length - 1]}`;
    }

    const auditPass = useCallback(() => {
        setAduitLoaing(true);
        //调用接口
        setAduitLoaing(false);
        props?.onCompleted?.();
    }, [props]);


    if (loading) {
        return <Spin spinning={true} tip='加载中,请稍后...' >
            <div style={{ height: 200 }} ></div>
        </Spin>
    }

    return <div className='auditPage' >
        <div style={{ display: 'flex' }} >
            <div>
                <Image width={100} src={pngLogo} />
            </div>
            <div style={{ marginRight: 50 }} >
                <div className='appDesc' >
                    <div className='appDescLbl' >名称：</div>
                    <div>这是应用名称</div>
                </div>
                <div className='appDesc' >
                    <div className='appDescLbl' >版本：</div>
                    <div>0.0.1</div>
                </div>
                <div className='appDesc' >
                    <div className='appDescLbl' >短描述：</div>
                    <div>短描述最多九个字捏</div>
                </div>
            </div>
            <div>
                <div style={{ marginBottom: 12 }} >
                    <Button type='primary' danger onClick={() => refuseRef.current?.show(appInfo)} >审核拒绝</Button>
                </div>
                <div  >
                    <Button type='primary' loading={aduitLoading} onClick={() => auditPass()}  >审核通过</Button>
                </div>
            </div>
        </div>
        <div>
            <Tabs
                activeKey={activeTabKey}
                onChange={(e: any) => setActiveTabKey(e)}
                items={getTabItem()}
            />
        </div>
        <div>
            {
                activeTabKey === TabKey.base &&
                <div>
                    <div style={{ marginBottom: 12 }} >
                        应用类别：
                        这是类别
                        &nbsp;&nbsp;
                        <Tooltip title='没有这个类别捏' >
                            <span style={{ color: 'red', cursor: 'pointer' }} ><ExclamationCircleOutlined /></span>
                        </Tooltip>
                    </div>
                    <div>
                    </div>
                    <div>
                        应用连接：
                        <a>https://www.baidu.com</a>
                    </div>
                </div>
            }
            {
                activeTabKey === TabKey.info &&
                <div>
                    <div >
                        {
                            (testPicList).map((m) => <span style={{ display: 'inline-block', marginLeft: 12, width: 100, height: 100, padding: 8, border: '1px solid #cecece' }} >
                                <div style={{ height: '100%', display: 'flex', justifyContent: 'center', alignItems: 'center' }} > <Image src={m.src} /></div>
                            </span>)
                        }
                    </div>
                    <div style={{ marginTop: 23, marginBottom: 20 }} dangerouslySetInnerHTML={{ __html: '这里是内容内容内容' }} ></div>
                </div>
            }
            {
                activeTabKey === TabKey.contributors &&
                <div>
                    {
                        (testContributors).map((m) => <div className='appCon' >
                            <div>
                                <Avatar style={{ background: '#1890ff' }} >{getUserNameLastDigit(m.userName)}</Avatar>
                            </div>
                            <div style={{ marginLeft: 10, marginTop: -3 }} >
                                <div><span style={{ fontWeight: 'bolder' }} >{m.userNick}</span>({m.userName})</div>
                                <div style={{ color: '#8c8c8c', fontSize: 12 }} >{m.roleName}</div>
                            </div>
                        </div>)
                    }
                </div>
            }
            {
                activeTabKey === TabKey.version &&
                <div>
                    {
                        (testVersions).map((m) => <div className='appVersionRec' >
                            <div className='appVersionRecTitle' >{m.title}</div>
                            <div>{m.describe}</div>
                        </div>)
                    }
                </div>
            }
        </div>
        <AuditRefuseModal ref={refuseRef} onCompleted={() => props.onCompleted?.()} />
    </div>
}

export default AuditPage;

export interface IInsertKeyTipProps {
    /**
     * 操作确认后的回调
     */
    onConfirm?: (data: IInsertKeyTipConfirmData) => void;
    /**
     * 操作取消后的回调
     */
    onCancel?: () => void;
    /**
     * 提示模式 仅提示插key 仅提示密码 两者都展示
     */
    mode?: 'onlyTip' | 'onlyPwd' | 'both';
    /**
     * 获取设备号的方法
     */
    deviceGet?: () => Promise<string>;
    /**
     * 验证密码的方法
     */
    pwdCheck?: (pwd: string) => Promise<boolean>;
    /**
     * 插入key提示标题，默认为 请将印章介质插入usb口
     */
    insertTitle?: string;
    /**
     * 当插入设备有误时，提示内容，默认为 请插入正确的介质
     */
    devError?: string;
    /**
     * 密码框title 默认为 输入介质密码
     */
    pwdTitle?: string;
    /**
     * 当输入密码有误时，提示内容，默认为 介质密码错误
     */
    pwdError?: string;
}

export interface IInsertKeyTipConfirmData {
    /**
     * 设备号 需show时传入，或者提供了获取设备号方法，才会提供
     */
    deviceId?: string;
    /**
     * 密码
     */
    password?: string;
    /**
     * 调用show方法时传入的其他数据对象，传入什么这里就吐出什么
     */
    other?: any;
}
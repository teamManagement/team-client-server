import React from "react"


export interface IFileUploadProps {
    /**
     * 值
     */
    value?: string | ArrayBuffer | File | IFileUploadFullRsp | null,
    /**
     * onChange 事件
     */
    onChange?: (value: string | ArrayBuffer | File | IFileUploadFullRsp | null) => void,
    /**
     * value值的模式 默认为base64
     */
    valueMode?: 'base64' | 'arrayBuffer' | 'file' | 'full-base64' | 'full-arrayBuffer',
    /**
     * 选择文件后文件信息展示模式，默认为 other
     */
    showMode?: 'picture' | 'other',
    /**
     * picture 模式下是否可以预览图片，默认为true
     */
    previewPic?: boolean,
    /**
     * 文件最大体积，单位 字节
     */
    maxSize?: number,
    /**
     * 超过最大体积限制后的提示内容
     */
    overSizeTip?: React.ReactNode,
    /**
     * 指定展示文件名(在没有文件名却又要展示文件名的时候显示)
     */
    fileName?: string,
    /**
     * 提示标头 默认为 ‘选择文件’
     */
    tipTitle?: string,
    /**
     * 提示内容，常用来提示文件格式，大小
     */
    placeHolder?: string,
    /**
     * other模式下 选择文件后，预览区域内容 不传默认为文件图标
     */
    previewIcon?: React.ReactNode,
    /**
     * 可选的文件类型
     */
    accept?: string,
    /**
     * 是否禁用
     */
    disabled?: boolean,
}

export interface IFileUploadFullRsp {
    /**
     * 真实文件名
     */
    name: string,
    /**
     * props里指定的文件名
     */
    orderName: string | undefined,
    /**
     * 后缀名 不带点
     */
    suffixName: string,
    /**
     * HTTP content-type
     */
    contentType: string,
    /**
     * 文件大小
     */
    size: number,
    /**
     * 文件对象
     */
    source: string | ArrayBuffer
}
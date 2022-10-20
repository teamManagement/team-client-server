import React, { Component } from 'react';
import { Modal, Image, Tooltip } from 'antd';
import { CloseCircleFilled, CloudUploadOutlined } from '@ant-design/icons';
import { IFileUploadFullRsp, IFileUploadProps } from '../types';
import { newGuid, isNull, isNullArray, getBase64SizeSync } from '../utils';
import { FileOutlined, LoadingOutlined } from '@ant-design/icons'
import './index.less';


interface FileUploadState_FileObj {
    name?: string,
    size?: number,
    valueSource?: any
}

interface IState {
    inputId: string,
    fileObj: FileUploadState_FileObj,
}

export default class FileUpload extends Component<IFileUploadProps, IState> {
    constructor(props: IFileUploadProps) {
        super(props);
        this.state = {
            inputId: newGuid(),
            fileObj: {},
        }
    }
    getInput(): HTMLInputElement {
        return (document.getElementById(this.state.inputId) as HTMLInputElement);
    }
    clearFile() {
        this.getInput().value = '';
        this.setState({ fileObj: {} });
        this.props.onChange?.(null);
    }
    onFileChange(files: FileList) {
        const { valueMode, onChange, maxSize, overSizeTip } = this.props;
        if (isNullArray(files)) {
            this.clearFile();
            return;
        }
        const file = files[0];
        if ((maxSize && maxSize > 0) && file.size > maxSize) {
            this.clearFile();
            Modal.error({ title: overSizeTip ? overSizeTip : '文件大小超过最大限制' });
            return;
        }
        var fileObj: FileUploadState_FileObj = {
            name: file.name,
            size: file.size
        }
        //初步更新
        this.setState({ fileObj });
        var valueFlag = valueMode ?? 'base64';
        if (valueFlag === 'file') {
            onChange?.(file);
            fileObj.valueSource = file;
            this.setState({ fileObj });
            return;
        }
        const reader = new FileReader();
        reader.onload = (ev) => {
            const readRsp: any = ev.target?.result;
            if (valueFlag === 'base64' || valueFlag === 'arrayBuffer') {
                onChange?.(readRsp);
                fileObj.valueSource = readRsp;
                this.setState({ fileObj });
                return;
            }
            var fullRsp: IFileUploadFullRsp = {
                name: file.name,
                orderName: this.props.fileName,
                suffixName: file.name.slice(file.name.lastIndexOf('.') + 1),
                contentType: file.type,
                size: file.size,
                source: readRsp
            };
            onChange?.(fullRsp);
            fileObj.valueSource = fullRsp;
            this.setState({ fileObj });
        };
        if (valueFlag === 'base64' || valueFlag === 'full-base64') {
            reader.readAsDataURL(file);
        } else if (valueFlag === 'arrayBuffer' || valueFlag === 'full-arrayBuffer') {
            reader.readAsArrayBuffer(file);
        }
    }
    picLoading: boolean = false;
    getImgSrc(value: any): any {
        const valueSource = value ?? this.state.fileObj?.valueSource;
        const { showMode, valueMode } = this.props;
        if (showMode !== 'picture') {
            return;
        }
        if (this.picLoading) {
            return;
        }
        this.picLoading = true;
        var rst = undefined;
        const valueFlag = valueMode ?? 'base64';
        if (isNull(valueSource)) {
            rst = undefined;
        }
        else if (valueFlag === 'base64') {
            rst = valueSource;
        }
        else if (valueFlag === 'arrayBuffer') {
            rst = window.URL.createObjectURL(new Blob([valueSource]));
        }
        else if (valueFlag === 'file') {
            rst = window.URL.createObjectURL(valueSource);
        }
        else if (valueFlag === 'full-base64') {
            rst = valueSource.source;
        }
        else if (valueFlag === 'full-arrayBuffer') {
            rst = window.URL.createObjectURL(new Blob([valueSource.source]));
        }
        this.picLoading = false;
        return rst;
    }
    sizeReading: boolean = false;
    getFileSize(value: any) {
        const valueSource = value ?? this.state.fileObj?.valueSource;
        if (this.sizeReading) {
            return;
        }
        this.sizeReading = true;
        const { fileObj } = this.state;
        const { valueMode } = this.props;
        var size = null;
        const valueFlag = valueMode ?? 'base64';
        if (isNull(valueSource)) {
            size = null;
        }
        else if (valueFlag === 'base64') {
            size = getBase64SizeSync(valueSource, true);
        }
        else if (valueFlag === 'arrayBuffer') {
            size = valueSource.byteLength;
        }
        else if (valueFlag === 'file') {
            size = valueSource.size;
        }
        else if (valueFlag === 'full-base64') {
            size = getBase64SizeSync(valueSource.source, true);
        }
        else if (valueFlag === 'full-arrayBuffer') {
            size = valueSource.source.byteLength;
        }
        this.sizeReading = false;
        return size;
    }
    render() {
        const { accept, tipTitle, previewIcon, placeHolder, disabled, value, showMode, previewPic, fileName } = this.props;
        const { inputId, fileObj } = this.state;
        const hasFile = !isNull(value) || !isNull(fileObj);
        const showModeFlag = showMode ?? 'other';
        const fsName = !isNull(fileObj?.name) ? fileObj?.name : fileName;
        let fileSize = null;
        let imgSrc = null;
        if (hasFile) {
            imgSrc = this.getImgSrc(value);
            fileSize = this.getFileSize(value);
        }

        return (
            <div className='byutils-fileUpload' >
                <input style={{ display: 'none' }} type='file' id={inputId} accept={accept} onChange={(e: any) => this.onFileChange(e.target.files)} />
                <div className='byutils-fileUpload-uploadBox' onClick={() => !disabled && this.getInput()?.click()} style={disabled ? { cursor: 'not-allowed' } : {}} >
                    <div className='byutils-fileUpload-uploadIconBox' >
                        <div className='byutils-fileUpload-uploadIcon' ><CloudUploadOutlined /></div>
                    </div>
                    <div className='byutils-fileUpload-uploadTipBox' >
                        <div className='byutils-fileUpload-uploadTipTitle' >{hasFile ? '重新上传' : !isNull(tipTitle) ? tipTitle : '打开文件'}</div>
                        <div className='byutils-fileUpload-uploadTipContent' >{placeHolder}</div>
                    </div>
                </div>
                {
                    hasFile &&
                    <>
                        <div className='byutils-fileUpload-previewBox' >
                            <div className='byutils-fileUpload-previewBody' >
                                {
                                    showModeFlag === 'other' && <span style={{ fontSize: 52 }} >{previewIcon ?? <FileOutlined />}</span>
                                }
                                {
                                    showModeFlag === 'picture' && (
                                        isNull(imgSrc) ? <span style={{ fontSize: 52 }} ><LoadingOutlined /></span>
                                            :
                                            <Image src={imgSrc} preview={isNull(previewPic) ? true : previewPic} />
                                    )
                                }
                            </div>
                            <div className='byutils-fileUpload-previewDesc' >
                                {/* <div className='upload-file'>{'打开'}</div> */}
                                <div className='byutils-fileUpload-previewDesc-fileName' >
                                    <Tooltip placement='top' >
                                        {fsName}
                                    </Tooltip>
                                </div>
                                <div className='byutils-fileUpload-previewDesc-fileSize' >
                                    {
                                        isNull(fileSize)
                                            ?
                                            <span><LoadingOutlined />&nbsp;</span>
                                            :
                                            (this.getFileSize(value) / 1024)?.toFixed(2)
                                    } KB
                                </div>
                            </div>
                        </div>
                        {
                            !disabled &&
                            <div className='byutils-fileUpload-cancel' onClick={() => this.clearFile()} >
                                <CloseCircleFilled />
                            </div>
                        }
                    </>
                }
            </div>
        )
    }
}
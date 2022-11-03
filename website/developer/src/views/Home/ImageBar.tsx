import { PlusOutlined } from '@ant-design/icons';
import { Image, Modal, Upload } from 'antd';
import type { RcFile, UploadProps } from 'antd/es/upload';
import type { UploadFile } from 'antd/es/upload/interface';
import { FC, useCallback, useEffect, useState } from 'react';

const getBase64 = (file: RcFile): Promise<string> =>
  new Promise((resolve, reject) => {
    const reader = new FileReader();
    reader.readAsDataURL(file);
    reader.onload = () => resolve(reader.result as string);
    reader.onerror = error => reject(error);
  });

interface IImageProps {
  type: 'single' | 'wall',
  getFileMenuList: () => void,
  getId: any,
}

const ImageBarWar: FC<IImageProps> = (props) => {

  const [previewOpen, setPreviewOpen] = useState(false);
  const [previewImage, setPreviewImage] = useState('');
  const [previewTitle, setPreviewTitle] = useState('');
  const [fileList, setFileList] = useState<UploadFile[]>([]);
  const [saveFile, setSaveFile] = useState<any[]>([])


  const handleCancel = () => setPreviewOpen(false);

  const handlePreview = async (file: UploadFile) => {
    if (!file.url && !file.preview) {
      file.preview = await getBase64(file.originFileObj as RcFile);
    }

    setPreviewImage(file.url || (file.preview as string));
    setPreviewOpen(true);
    setPreviewTitle(file.name || file.url!.substring(file.url!.lastIndexOf('/') + 1));
  };

  const handleChange: UploadProps['onChange'] = ({ fileList: newFileList }) =>
    setFileList(newFileList);


  const uploadButton = (
    <div>
      <PlusOutlined />
      <div style={{ marginTop: 8 }}>上传图片</div>
    </div>
  );
  const onFileChange = useCallback(async (fileInfo: any) => {
    console.log(fileInfo);
    if (props.type === 'single') {
      const appInfo: any = await window.teamworkSDK.store.get(props.getId)
      const imageId = await window.teamworkSDK.cache.file.store(fileInfo?.path, -1)
      const getImage = window.teamworkSDK.cache.file.getDownloadUrl(imageId)
      const newAppInfo = { ...appInfo, icon: getImage }
      await window.teamworkSDK.store.set(props.getId, newAppInfo)
      const appInfos: any = await window.teamworkSDK.store.get(props.getId)
      const allList: any = await window.teamworkSDK.store.get('_content_menu_list')
      const list = allList.map((m: any) => m.id === appInfos.id ? appInfos : m)
      await window.teamworkSDK.store.set('_content_menu_list', list)
      setPreviewImage(getImage)
      await window.teamworkSDK.store.set('image-single', getImage)
      props.getFileMenuList()
    }
    // if (props.type === 'single') {
    // const imageId = await window.teamworkSDK.cache.file.store(fileInfo?.path, -1)
    // const getImage = window.teamworkSDK.cache.file.getDownloadUrl(imageId)
    // setPreviewImage(getImage)
    // window.teamworkSDK.store.set('image-single', getImage)
    //   window.location.reload()
    // } else {
    //   let newFile = [...saveFile]
    //   const imageId = await window.teamworkSDK.cache.file.store(fileInfo?.path, -1)
    //   const getImage = window.teamworkSDK.cache.file.getDownloadUrl(imageId)
    //   setPreviewImage(getImage)
    //   newFile.push(getImage)
    //   setSaveFile(newFile)
    //   window.teamworkSDK.store.set('image', newFile)
    // }
  }, [saveFile, props])


  const getImageFirst = useCallback(async () => {

    if (props.type === 'single') {
      const appInfo: any = await window.teamworkSDK.store.get(props.getId)
      setFileList([
        {
          uid: '-1',
          name: 'image.png',
          status: 'done',
          url: appInfo?.icon,
        }
      ])
    }
    // if (props.type === 'single') {
    //   const list: any = await window.teamworkSDK.store.get<{ [key: string]: string }>('image-single')
    //   console.log(list);
    //   if (list) {
    //     setFileList([
    //       {
    //         uid: '-1',
    //         name: 'image.png',
    //         status: 'done',
    //         url: list,
    //       }
    //     ])
    //   }

    // } else {
    //   const list: any = await window.teamworkSDK.store.get<{ [key: string]: string }>('image')
    //   console.log(list);

    //   if (list) {
    //     setPreviewImage(list)
    //     setFileList(
    //       list?.map((m: any) => {
    //         return {
    //           uid: '-1',
    //           name: 'image.png',
    //           status: 'done',
    //           url: m,
    //         }
    //       })
    //     )
    //   }
    // }
  }, [props])

  useEffect(() => {
    getImageFirst()
  }, [getImageFirst])

  return (
    <>
      <Upload
        listType="picture-card"
        fileList={fileList}
        beforeUpload={(file) => onFileChange(file)}
        onPreview={handlePreview}
        onChange={handleChange}
      >
        {props.type === 'single' ? (fileList.length >= 1 ? null : uploadButton)
          : (fileList.length >= 9 ? null : uploadButton)
        }
      </Upload>
      <Modal open={previewOpen} title={previewTitle} footer={null} onCancel={handleCancel}>
        <Image height={120} alt="example" preview={false} style={{ width: '100%' }} src={previewImage} />
      </Modal>
    </>
  )
}

export default ImageBarWar
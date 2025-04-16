import React, { useState } from 'react';
import { filesAPI } from '../../api/api';
import { useRouteLoaderData } from 'react-router-dom';

const FileAttachModal = ({ show, handleClose, entity, entityUuid }) => {
    const [selectedFile, setSelectedFile] = useState(null);
  const [title, setTitle] = useState('');

  const handleFileChange = (e) => {
    setSelectedFile(e.target.files[0]);
  };

  const handleUpload = async () => {
    if (!selectedFile) return alert("Выберите файл");

    const formData = new FormData();
    formData.append('file', selectedFile);


    try {
      const response = await filesAPI.uploadFile(formData);
      if (response.status === 200) {
        const fileUrl = response.data.file_url;
        const attachResponse = await filesAPI.addFile(
            entity,
            entityUuid,
            {   
                title: title,
                file_url: fileUrl 
            }
        );

        if (attachResponse.status === 200) {
            alert('Файл успешно загружен и прикреплен!');
            handleClose();
          } else {
            alert('Ошибка при прикреплении файла.');
          }
        }
        else {
            alert('Ошибка при прикреплении файла.');
        }
        
        } catch (error) {
          console.error('Ошибка при загрузке файла:', error);
          alert('Ошибка при загрузке файла.');
        }
  };

  if (!show) return null;

  return (
    <div className="modal-overlay" onClick={handleClose}>
      <div className="modal-content" onClick={e => e.stopPropagation()}>
        <span className="close" onClick={handleClose}>&times;</span>
        <h2>Прикрепить файл</h2>
        <input type="text" placeholder="Название файла" value={title} onChange={e => setTitle(e.target.value)} />
        <input type="file" onChange={handleFileChange} />
        <button onClick={handleUpload}>Загрузить</button>
        <button onClick={handleClose}>Отмена</button>
      </div>
    </div>
  );
};

export default FileAttachModal


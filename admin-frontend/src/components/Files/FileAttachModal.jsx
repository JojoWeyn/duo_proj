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
    <div className="modal-overlay open" onClick={handleClose}>
      <div className="modal-content" onClick={e => e.stopPropagation()}>
        <div className="modal-header">
          <h2 className="modal-title">Прикрепить файл</h2>
          <span className="modal-close" onClick={handleClose}>&times;</span>
        </div>
        <div className="modal-form">
          <input type="text" placeholder="Название файла" value={title} onChange={e => setTitle(e.target.value)} />
          <input type="file" onChange={handleFileChange} />
          <div className="modal-actions">
            <button onClick={handleClose}>Отмена</button>
            <button onClick={handleUpload} className="primary">Загрузить</button>
          </div>
        </div>
      </div>
    </div>
  );
};

export default FileAttachModal


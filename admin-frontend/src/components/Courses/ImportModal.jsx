import React, { useState } from 'react';
import { coursesAPI } from '../../api/api';

const ImportExcelModal = ({ show, handleClose }) => {
  const [selectedFile, setSelectedFile] = useState(null);

  const handleFileChange = (e) => {
    setSelectedFile(e.target.files[0]);
  };

  const handleUpload = async () => {
    if (!selectedFile) return alert("Выберите файл");

    const formData = new FormData();
    formData.append('file', selectedFile);

    try {
      const response = await coursesAPI.importExcel(formData);
      if (response.status === 200) {
        alert('Импорт успешно выполнен!\nКурс появится в течении минуты');
        handleClose();
      } else {
        alert('Ошибка при импорте файла.');
      }
    } catch (error) {
      console.error('Ошибка при импорте:', error);
      alert('Ошибка при импорте файла.');
    }
  };

  if (!show) return null;

  return (
    <div className="modal-overlay open" onClick={handleClose}>
      <div className="modal-content" onClick={(e) => e.stopPropagation()}>
        <div className="modal-header">
          <h2 className="modal-title">Импортировать Excel</h2>
          <span className="modal-close" onClick={handleClose}>&times;</span>
        </div>

        <div className="modal-form">
          <input type="file" accept=".xlsx,.xls" onChange={handleFileChange} />

          <a
            href="https://s3.regru.cloud/duo/Книга1.xlsx"
            download
            className="modal-sample-link"
            style={{ marginTop: '10px', color: '#007bff', textDecoration: 'underline' }}
          >
            📥 Скачать образец Excel-файла
          </a>

          <div className="modal-actions">
            <button onClick={handleClose}>Отмена</button>
            <button onClick={handleUpload} className="primary">Загрузить</button>
          </div>
        </div>
      </div>
    </div>
  );
};

export default ImportExcelModal;
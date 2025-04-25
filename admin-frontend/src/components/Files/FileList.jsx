import React, { useEffect, useState } from 'react';
import { filesAPI } from '../../api/api';
import './FileList.css';

const getFileType = (url) => {
  const ext = url.split('.').pop().split('?')[0].toLowerCase();
  if (['jpg', 'jpeg', 'png', 'gif', 'webp'].includes(ext)) return 'image';
  if (['mp4', 'webm', 'mov'].includes(ext)) return 'video';
  if (['pdf', 'doc', 'docx', 'txt', 'xlsx'].includes(ext)) return 'document';
  return 'other';
};

export const FileList = () => {
  const [files, setFiles] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    const fetchFiles = async () => {
      try {
        setLoading(true);
        setError(null);
        const response = await filesAPI.getList();
        const data = await response.data;
        setFiles(data.files);
      } catch (err) {
        console.error('Ошибка загрузки файлов:', err);
        setError('Не удалось загрузить файлы.');
      } finally {
        setLoading(false);
      }
    };

    fetchFiles();
  }, []);

  const handleDeleteClick = async (url) => {
    const fileName = url.substring(url.indexOf('duo/') + 4);
    const isConfirmed = window.confirm("Вы уверены, что хотите удалить этот файл?");
    if (!isConfirmed) return;

    try {
      await filesAPI.deleteFile(fileName);
      setFiles(files.filter((file) => file !== url));
    } catch (error) {
      console.error("Ошибка при удалении файла:", error);
    }
  };

  return (
    <div className="file-list-container">
      {loading && <p>Загрузка файлов...</p>}
      {error && <div className="error">{error}</div>}
      {!loading && !error && files.length === 0 && <p>Файлы не найдены</p>}

      {!loading && !error && files.map((url, idx) => {
        const type = getFileType(url);
        return (
          <div key={idx}>
            <div className="file-card">
              {type === 'image' && (
                <img src={url} alt={`file-${idx}`} className="file-image" />
              )}
              {type === 'video' && (
                <video controls className="file-video">
                  <source src={url} />
                  Ваш браузер не поддерживает видео.
                </video>
              )}
              {type === 'document' && (
                <a
                  href={url}
                  target="_blank"
                  rel="noreferrer"
                  className="file-link document-link"
                >
                  📄 Открыть документ
                </a>
              )}
              {type === 'other' && (
                <a
                  href={url}
                  target="_blank"
                  rel="noreferrer"
                  className="file-link other-link"
                >
                  📁 {url.split('/').pop()}
                </a>
              )}
            </div>
            <button className="delete-button full-width" onClick={() => handleDeleteClick(url)}>
              Удалить
            </button>
          </div>
        );
      })}
    </div>
  );
};
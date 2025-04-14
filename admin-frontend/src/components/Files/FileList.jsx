import React, { useEffect, useState } from 'react';
import { filesAPI } from '../../api/api';

const getFileType = (url) => {
  const ext = url.split('.').pop().split('?')[0].toLowerCase();
  if (['jpg', 'jpeg', 'png', 'gif', 'webp'].includes(ext)) return 'image';
  if (['mp4', 'webm', 'mov'].includes(ext)) return 'video';
  if (['pdf', 'doc', 'docx', 'txt', 'xlsx'].includes(ext)) return 'document';
  return 'other';
};

export const FileList = () => {
  const [files, setFiles] = useState([]);
  const [error, setError] = useState('');

  useEffect(() => {
    const fetchFiles = async () => {
      try {
        const response = await filesAPI.getList();
        const data = await response.data;
        setFiles(data.files);
      } catch (err) {
        console.error('Ошибка загрузки файлов:', err);
        setError('Не удалось загрузить файлы.');
      }
    };

    fetchFiles();
  }, []);

  if (error) return <div className="error">{error}</div>;

  return (
    <div className="file-list-container">
      {files.map((url, idx) => {
        const type = getFileType(url);

        return (
          <div key={idx} className="file-card">
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
        );
      })}
    </div>
  );
};

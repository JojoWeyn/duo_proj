import React from 'react';

const ConfirmDeleteModal = ({ show, onConfirm, onCancel }) => {
  if (!show) return null;

  return (
    <div className="modal-overlay open" onClick={onCancel}>
      <div className="modal-content small" onClick={(e) => e.stopPropagation()}>
        <div className="modal-header">
          <h2 className="modal-title">Удалить?</h2>
          <span className="modal-close" onClick={onCancel}>&times;</span>
        </div>
        <div className="modal-body">
          <p>Вы уверены, что хотите удалить? <br></br>Это действие необратимо.</p>
        </div>
        <div className="modal-actions">
          <button onClick={onCancel}>Нет</button>
          <button onClick={onConfirm} className="danger">Да</button>
        </div>
      </div>
    </div>
  );
};

export default ConfirmDeleteModal;
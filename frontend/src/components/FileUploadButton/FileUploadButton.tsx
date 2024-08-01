import { ChangeEvent, useState } from 'react';
import { UploadIcon } from '../Icons/Icons';
import yaml from 'js-yaml';
import './FileUploadButton.scss';

const FileUploadButton = () => {
  const [jsonData, setJsonData] = useState<Record<string, unknown> | null>(
    null
  );

  const handleFileUpload = (e: ChangeEvent<HTMLInputElement>) => {
    if (e.target.files && e.target.files.length > 0) {
      const file = e.target.files[0];

      if (
        file.type !== 'application/x-yaml' &&
        file.type !== 'application/yaml'
      )
        return;

      const reader = new FileReader();

      reader.onload = e => {
        try {
          if (!e.target) return;

          const yamlText = e.target.result as string;
          const jsonObject = yaml.load(yamlText) as Record<string, unknown>;

          setJsonData(jsonObject);
        } catch (err) {
          setJsonData(null);
        }
      };

      reader.readAsText(file);
    }
  };

  console.log('JSON', jsonData);

  return (
    <label className="file-upload-input-container" htmlFor="file-upload">
      <span className="file-upload-label">
        <UploadIcon height={22} width={22} /> Upload Workflow
      </span>
      <input
        accept=".yaml, .yml"
        className="file-upload-input"
        id="file-upload"
        onChange={handleFileUpload}
        type="file"
      />
    </label>
  );
};

export default FileUploadButton;

import imageCompression from "browser-image-compression";

const OPTIONS = {
  maxSizeMB: 1,
  maxWidthOrHeight: 2048,
  useWebWorker: true,
  initialQuality: 0.8,
};

export async function compressImage(file) {
  if (!file.type.startsWith("image/")) return file;
  if (file.size <= 200 * 1024) return file;
  try {
    return await imageCompression(file, OPTIONS);
  } catch {
    return file;
  }
}

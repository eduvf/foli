function icon(emoji) {
  const size = 32;

  let canvas = document.createElement('canvas');
  canvas.width = size;
  canvas.height = size;

  let context = canvas.getContext('2d');
  context.font = '22pt system-ui';
  context.textBaseline = 'middle';
  context.textAlign = 'center';
  context.fillText(emoji, size / 2, size / 2);

  let link = document.createElement('link');
  link.rel = 'icon';
  link.href = canvas.toDataURL();

  document.head.appendChild(link);
}
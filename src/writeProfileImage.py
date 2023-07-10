import os
import re
import base64
from sys import argv
from PIL import Image
from io import BytesIO

FILE_REMOVE = None


def main() -> None:
    try:
        avatar_id, path = get_folder()
        text = re.sub('^data:image/.+;base64,', '', getBase64Text())
        img = Image.open(BytesIO(base64.b64decode(text)))
        img = crop_center(img)
        for size, format in ((650, 'webp'), (650, 'jpeg'), (200, 'webp'), (200, 'jpeg')):
            save_image(path, img.copy(), size, format)
        print(avatar_id)
    except Exception as err:
        print(f'ERROR: "{err}"')
        os.rmdir(path)
        exit()
    FILE_REMOVE()


def getBase64Text() -> str:
    """Получение текста base64 изображения из файла."""
    file_name = os.path.join('src', 'static', 'images',
                             'profile', argv[1])
    with open(file_name, 'r') as file:
        text = file.read()
    if text == 'none':
        os.remove(file_name)
        print('ERROR: "Text of None"')
        exit()
    global FILE_REMOVE
    def FILE_REMOVE(): return os.remove(file_name)
    return text


def save_image(path, img, size, format) -> None:
    """Cохраняет изображения в файлы разных форматов"""
    file_name = str(size) + '.' + (format if format != 'jpeg' else 'jpg')
    img.thumbnail(size=(size, size))
    img.save(os.path.join(path, file_name), format=format, quality=95)


def get_folder() -> tuple[int, str]:
    """Возвращает id аватара и директорию с аватаром"""
    main_path = os.path.join('src', 'static', 'images', 'profile', '')
    max_folder = 1
    for i in os.listdir(main_path):
        path = main_path + i
        if os.path.isdir(path):
            if i.isnumeric():
                i = int(i)
                if i > max_folder:
                    max_folder = i

    if os.path.isdir(main_path + str(max_folder)):
        max_folder += 1
    os.mkdir(main_path + str(max_folder))
    return max_folder, main_path + str(max_folder)


def crop_center(pil_img) -> Image:
    """Функция для обрезки изображения по центру"""
    crop = min(pil_img.size)
    img_width, img_height = pil_img.size
    return pil_img.crop(((img_width - crop) // 2,
                        (img_height - crop) // 2,
                        (img_width + crop) // 2,
                        (img_height + crop) // 2))


if __name__ == '__main__':
    main()

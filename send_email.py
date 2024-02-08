from application.mail import Email


print('[CucinAssistant] Invio email')

okay = False
while not okay:
    print('\nInserisci il titolo della mail: ')
    title = input('> ')
    okay = input('Corretto? [y/n] ') == 'y'

okay = False
while not okay:
    print('\nInserisci il contenuto della mail: ')
    content = input('> ')

    print(f'\n---\n\nContenuto della mail:\n{content}\n')
    okay = input('Corretto? [y/n] ') == 'y'

print(title, content)
print('\nInvio della mail...', end=' ')
n = Email(title, 'base', content=content).broadcast()
print(f'Fatto! {n} email inviate.')

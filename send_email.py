from application.mail import BroadcastEmail


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
n = BroadcastEmail(title, content).send_all()
print(f'Fatto! {n} email inviate.')

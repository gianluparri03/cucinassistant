from cucinassistant.email import Email
from cucinassistant.config import config


print('[CucinAssistant] Invio email broadcast')

print('\nInserisci il titolo della mail: ')
title = input('> ')
print('\nInserisci il corpo della mail: ')
content = input('> ')

mail = Email(title, 'base', content=content)

print('\nInvio della mail di prova...', end=' ')
mail.send(config['Email']['Address'])
print('fatto.')

print("\nProseguire con l'invio in broadcast? [prosegui]")

if input('> ') == 'prosegui':
    n = mail.broadcast()
    print(f'\nFatto! {n} email inviate.')
else:
    print('\nInvio cancellato.')

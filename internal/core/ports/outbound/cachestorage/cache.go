package cachestorage

import (
	"context"
	"time"
)

type CacheStorage interface {
	Close() error

	// Get recupera o valor associado a uma chave no sistema de cache.
	//
	// Parâmetros:
	// - ctx: Contexto para controle de timeout/cancelamento.
	// - key: A chave cujo valor será recuperado.
	//
	// Retorno:
	// - string: O valor armazenado na chave. Retorna uma string vazia ("") se a chave não for encontrada.
	// - error: Um erro é retornado em caso de falha na comunicação com o backend de cache ou erro interno.
	Get(ctx context.Context, key string) (string, error)

	// Set define um valor no cache com uma determinada chave e tempo de expiração.
	//
	// Parâmetros:
	// - ctx: Contexto para controle de timeout/cancelamento.
	// - key: A chave que será usada para armazenar o valor.
	// - value: O valor a ser armazenado (qualquer tipo serializável).
	// - expiration: Duração após a qual a chave expirará. Use 0 para chave sem expiração.
	//
	// Retorno:
	// - string: Retorna o status da operação, quando disponível.
	// - error: Um erro é retornado em caso de falha na operação.
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) (string, error)

	// SetNX define um valor no cache somente se a chave ainda não existir.
	//
	// É útil para implementar locks ou garantir que valores não sejam sobrescritos.
	//
	// Parâmetros:
	// - ctx: Contexto para controle de timeout/cancelamento.
	// - key: A chave a ser definida.
	// - value: O valor a ser armazenado.
	// - expiration: Tempo de expiração da chave.
	//
	// Retorno:
	// - bool: `true` se a chave foi definida com sucesso, `false` se já existia.
	// - error: Um erro é retornado em caso de falha na operação.
	SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error)

	// Incr incrementa o valor inteiro armazenado na chave especificada.
	//
	// Se a chave não existir, ela será criada com o valor inicial 0 antes de ser incrementada.
	// Este comando é atômico no backend de cache, sendo seguro para uso em ambientes concorrentes.
	//
	// Parâmetros:
	// - ctx: Contexto para controle de timeout/cancelamento.
	// - key: A chave cujo valor será incrementado.
	//
	// Retorno:
	// - int64: Novo valor da chave após o incremento.
	// - error: Um erro é retornado em caso de falha na operação ou se o valor atual da chave não for numérico.
	Incr(ctx context.Context, key string) (int64, error)

	// Expire define um tempo de expiração (TTL) para a chave especificada.
	//
	// Após o tempo definido, a chave será automaticamente removida do cache.
	// Útil para definir validade de dados temporários como caches, contadores e flags de jobs.
	//
	// Parâmetros:
	// - ctx: Contexto para controle de timeout/cancelamento.
	// - key: A chave cujo TTL será definido.
	// - ttl: Duração do tempo de vida da chave.
	//
	// Retorno:
	// - error: Um erro é retornado em caso de falha na operação de expiração.
	Expire(ctx context.Context, key string, ttl time.Duration) error

	// WaitForCacheValue realiza polling no sistema de cache para aguardar que uma chave atenda a uma condição específica.
	//
	// Essa função verifica periodicamente se o valor associado à chave `key` está disponível
	// e satisfaz a condição definida pela função `predicate`. O polling é feito em intervalos regulares
	// definidos por `interval`, e será encerrado se o valor desejado for encontrado ou se o tempo total
	// exceder o `timeout`.
	//
	// Parâmetros:
	// - ctx: Contexto que pode ser usado para cancelar manualmente a operação.
	// - key: A chave cujo valor será verificado.
	// - interval: Intervalo entre as tentativas de leitura do valor da chave.
	// - timeout: Tempo máximo de espera até a função retornar erro por timeout.
	// - predicate: Função que avalia se o valor lido do cache satisfaz a condição desejada.
	//
	// Retornos:
	// - string: O valor da chave que satisfaz a condição.
	// - error: Um erro será retornado em caso de timeout, falha de leitura ou erro no predicate.
	//
	// Exemplo de uso:
	//   val, err := cache.WaitForCacheValue(ctx, "job:status:123", 500*time.Millisecond, 10*time.Second, func(v string) (bool, error) {
	//       return v == "completed", nil
	//   })
	//   if err != nil {
	//       log.Fatal(err)
	//   }
	//   fmt.Println("Job finalizado com status:", val)
	WaitForCacheValue(
		ctx context.Context,
		key string,
		interval time.Duration,
		timeout time.Duration,
		predicate func(string) (bool, error),
	) (string, error)
}
